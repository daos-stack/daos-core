/**
 * (C) Copyright 2016-2019 Intel Corporation.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *    http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 *
 * GOVERNMENT LICENSE RIGHTS-OPEN SOURCE SOFTWARE
 * The Government's rights to use, modify, reproduce, release, perform, display,
 * or disclose this software are subject to the terms of the Apache License as
 * provided in Contract No. B609815.
 * Any reproduction of computer software, computer software documentation, or
 * portions thereof marked with this legend must also reproduce the markings.
 */
/**
 * This file is part of daos_sr
 *
 * src/placement/pl_map.c
 */
#define D_LOGFAC	DD_FAC(placement)

#include "pl_map.h"
#include <gurt/hash.h>

extern struct pl_map_ops	ring_map_ops;

/** dictionary for all unknown placement maps */
struct pl_map_dict {
	/** type of the placement map */
	pl_map_type_t		 pd_type;
	/** customized functions */
	struct pl_map_ops	*pd_ops;
	/** name of the placement map */
	char			*pd_name;
};

/** array of defined placement maps */
static struct pl_map_dict pl_maps[] = {
	{
		.pd_type	= PL_TYPE_RING,
		.pd_ops		= &ring_map_ops,
		.pd_name	= "ring",
	},
	{
		.pd_type	= PL_TYPE_UNKNOWN,
		.pd_ops		= NULL,
		.pd_name	= "unknown",
	},
};


static int
pl_map_create_inited(struct pool_map *pool_map, struct pl_map_init_attr *mia,
		    struct pl_map **pl_mapp)
{
	struct pl_map_dict	*dict = pl_maps;
	struct pl_map		*map;
	int			 rc;

	for (dict = &pl_maps[0]; dict->pd_type != PL_TYPE_UNKNOWN; dict++) {
		if (dict->pd_type == mia->ia_type)
			break;
	}

	if (dict->pd_type == PL_TYPE_UNKNOWN) {
		D_DEBUG(DB_PL,
			"Unknown placement map type %d\n", dict->pd_type);
		return -EINVAL;
	}

	D_DEBUG(DB_PL, "Create a %s placement map\n", dict->pd_name);

	rc = dict->pd_ops->o_create(pool_map, mia, &map);
	if (rc != 0)
		return rc;

	rc = D_SPIN_INIT(&map->pl_lock, PTHREAD_PROCESS_PRIVATE);
	if (rc != 0) {
		dict->pd_ops->o_destroy(map);
		return rc;
	}

	map->pl_ref  = 1; /* for the caller */
	map->pl_connects = 0;
	map->pl_type = mia->ia_type;
	map->pl_ops  = dict->pd_ops;
	D_INIT_LIST_HEAD(&map->pl_link);

	*pl_mapp = map;
	return 0;
}

/**
 * Destroy a placement map
 */
void
pl_map_destroy(struct pl_map *map)
{
	D_ASSERT(map->pl_ref == 0);
	D_ASSERT(map->pl_ops != NULL);
	D_ASSERT(map->pl_ops->o_destroy != NULL);

	D_SPIN_DESTROY(&map->pl_lock);
	map->pl_ops->o_destroy(map);
}

/** Print a placement map, it's optional and for debug only */
void pl_map_print(struct pl_map *map)
{
	D_ASSERT(map->pl_ops != NULL);

	if (map->pl_ops->o_print != NULL)
		map->pl_ops->o_print(map);
}

/**
 * Compute layout for the input object metadata @md. It only generates the
 * layout of the redundancy group that @shard_md belongs to if @shard_md
 * is not NULL.
 */
int
pl_obj_place(struct pl_map *map, struct daos_obj_md *md,
	     struct daos_obj_shard_md *shard_md,
	     struct pl_obj_layout **layout_pp)
{
	D_ASSERT(map->pl_ops != NULL);
	D_ASSERT(map->pl_ops->o_obj_place != NULL);

	return map->pl_ops->o_obj_place(map, md, shard_md, layout_pp);
}

/**
 * Check if the provided object has any shard needs to be rebuilt for the
 * given rebuild version @rebuild_ver.
 *
 * \param  map [IN]		pl_map this check is performed on
 * \param  md  [IN]		object metadata
 * \param  shard_md [IN]	shard metadata (optional)
 * \param  rebuild_ver [IN]	current rebuild version
 * \param  tgt_rank [OUT]	spare target ranks
 * \param  shard_id [OUT]	shard ids to be rebuilt
 * \param  array_size [IN]	array size of tgt_rank & shard_id

 * \return	> 0	the array size of tgt_rank & shard_id, so it means
 *                      getting the spare targets for the failure shards.
 *		0	No need rebuild or find spare tgts successfully.
 *		-ve	error code.
 */
int
pl_obj_find_rebuild(struct pl_map *map, struct daos_obj_md *md,
		    struct daos_obj_shard_md *shard_md,
		    uint32_t rebuild_ver, uint32_t *tgt_rank,
		    uint32_t *shard_id, unsigned int array_size)
{
	D_ASSERT(map->pl_ops != NULL);

	if (!map->pl_ops->o_obj_find_rebuild)
		return -DER_NOSYS;

	return map->pl_ops->o_obj_find_rebuild(map, md, shard_md, rebuild_ver,
					       tgt_rank, shard_id, array_size);
}

/**
 * Check if the provided object shard needs to be built on the reintegrated
 * targets @tgp_reint.
 *
 * \return	1	Build the object on the returned target @tgt_reint.
 *		0	Skip this object.
 *		-ve	error code.
 */
int
pl_obj_find_reint(struct pl_map *map, struct daos_obj_md *md,
		  struct daos_obj_shard_md *shard_md,
		  struct pl_target_grp *tgp_reint, uint32_t *tgt_reint)
{
	D_ASSERT(map->pl_ops != NULL);

	if (!map->pl_ops->o_obj_find_reint)
		return -DER_NOSYS;

	return map->pl_ops->o_obj_find_reint(map, md, shard_md, tgp_reint,
					     tgt_reint);
}

void
pl_obj_layout_free(struct pl_obj_layout *layout)
{
	if (layout->ol_shards != NULL)
		D_FREE(layout->ol_shards);
	D_FREE(layout);
}

int
pl_obj_layout_alloc(unsigned int shard_nr, struct pl_obj_layout **layout_pp)
{
	struct pl_obj_layout *layout;

	D_ALLOC_PTR(layout);
	if (layout == NULL)
		return -DER_NOMEM;

	layout->ol_nr = shard_nr;
	D_ALLOC_ARRAY(layout->ol_shards, layout->ol_nr);
	if (layout->ol_shards == NULL)
		goto failed;

	*layout_pp = layout;
	return 0;
 failed:
	pl_obj_layout_free(layout);
	return -DER_NOMEM;
}

/**
 * Return the index of the first shard of the redundancy group that @shard
 * belongs to.
 */
unsigned int
pl_obj_shard2grp_head(struct daos_obj_shard_md *shard_md,
		      struct daos_oclass_attr *oc_attr)
{
	int sid	= shard_md->smd_id.id_shard;

	/* XXX: only for the static stripe classes for the time being */
	D_ASSERT(oc_attr->ca_schema == DAOS_OS_SINGLE ||
		 oc_attr->ca_schema == DAOS_OS_STRIPED);

	switch (oc_attr->ca_resil) {
	default:
		return sid;

	case DAOS_RES_EC:
	case DAOS_RES_REPL:
		return sid - sid % daos_oclass_grp_size(oc_attr);
	}
}
/**
 * Returns the redundancy group index of @shard_md.
 */
unsigned int
pl_obj_shard2grp_index(struct daos_obj_shard_md *shard_md,
		       struct daos_oclass_attr *oc_attr)
{
	int sid	= shard_md->smd_id.id_shard;

	/* XXX: only for the static stripe classes for the time being */
	D_ASSERT(oc_attr->ca_schema == DAOS_OS_SINGLE ||
		 oc_attr->ca_schema == DAOS_OS_STRIPED);

	switch (oc_attr->ca_resil) {
	default:
		return sid; /* no protection */

	case DAOS_RES_EC:
	case DAOS_RES_REPL:
		return sid / daos_oclass_grp_size(oc_attr);
	}
}

/** serialize operations on pl_htable */
pthread_rwlock_t	pl_rwlock = PTHREAD_RWLOCK_INITIALIZER;
/** hash table for placement maps */
struct d_hash_table	pl_htable = {
	.ht_ops			  = NULL,
};

#define DSR_RING_DOMAIN		PO_COMP_TP_RACK

void
pl_map_attr_init(struct pool_map *po_map, pl_map_type_t type,
		 struct pl_map_init_attr *mia)
{
	memset(mia, 0, sizeof(*mia));

	switch (type) {
	default:
		D_ASSERTF(0, "Unknown placemet map type: %d.\n", type);
		break;

	case PL_TYPE_RING:
		mia->ia_type	     = PL_TYPE_RING;
		mia->ia_ring.domain  = DSR_RING_DOMAIN;
		mia->ia_ring.ring_nr = 1;
		break;
	}
}

struct pl_map *
pl_link2map(d_list_t *link)
{
	return container_of(link, struct pl_map, pl_link);
}

static unsigned int
pl_hop_key_hash(struct d_hash_table *htab, const void *key,
		unsigned int ksize)
{
	D_ASSERT(ksize == sizeof(uuid_t));
	return d_hash_string_u32((const char *)key, ksize);
}

static bool
pl_hop_key_cmp(struct d_hash_table *htab, d_list_t *link,
	       const void *key, unsigned int ksize)
{
	struct pl_map *map = pl_link2map(link);

	D_ASSERT(ksize == sizeof(uuid_t));
	return !uuid_compare(map->pl_uuid, key);
}

static void
pl_hop_rec_addref(struct d_hash_table *htab, d_list_t *link)
{
	struct pl_map *map = pl_link2map(link);

	D_SPIN_LOCK(&map->pl_lock);
	map->pl_ref++;
	D_SPIN_UNLOCK(&map->pl_lock);
}

static bool
pl_hop_rec_decref(struct d_hash_table *htab, d_list_t *link)
{
	struct pl_map	*map = pl_link2map(link);
	bool		 zombie;

	D_ASSERT(map->pl_ref > 0);

	D_SPIN_LOCK(&map->pl_lock);
	map->pl_ref--;
	zombie = (map->pl_ref == 0);
	D_SPIN_UNLOCK(&map->pl_lock);

	return zombie;
}

void
pl_hop_rec_free(struct d_hash_table *htab, d_list_t *link)
{
	struct pl_map *map = pl_link2map(link);

	D_ASSERT(map->pl_ref == 0);
	pl_map_destroy(map);
}

static d_hash_table_ops_t pl_hash_ops = {
	.hop_key_hash		= pl_hop_key_hash,
	.hop_key_cmp		= pl_hop_key_cmp,
	.hop_rec_addref		= pl_hop_rec_addref,
	.hop_rec_decref		= pl_hop_rec_decref,
	.hop_rec_free		= pl_hop_rec_free,
};

#define PL_HTABLE_BITS		7

static int
pl_htable_init()
{
	/* NB: this hash table is created on demand, it will never
	 * be destroyed.
	 */
	D_ASSERT(pl_htable.ht_ops == NULL);
	return d_hash_table_create_inplace(D_HASH_FT_NOLOCK,
					   PL_HTABLE_BITS, NULL,
					   &pl_hash_ops, &pl_htable);
}

/**
 * Create a placement map based on attributes in \a mia
 */
int
pl_map_create(struct pool_map *pool_map, struct pl_map_init_attr *mia,
	      struct pl_map **pl_mapp, bool hash)
{
	if (hash) {
		int rc = 0;

		D_RWLOCK_WRLOCK(&pl_rwlock);
		if (!pl_htable.ht_ops)
			rc = pl_htable_init();
		D_RWLOCK_UNLOCK(&pl_rwlock);

		if (rc != 0) {
			D_ERROR("pl_htable_init failed, rc %d.\n", rc);
			return rc;
		}
	}

	return pl_map_create_inited(pool_map, mia, pl_mapp);


}

/**
 * Generate a new placement map from the pool map @pool_map, and replace the
 * original placement map for the same pool.
 *
 * \param	uuid [IN]	uuid of \a pool_map
 * \param	pool_map [IN]	pool_map
 * \param	connect [IN]	from pool connect or not
 */
int
pl_map_update(uuid_t uuid, struct pool_map *pool_map, bool connect)
{
	d_list_t		*link;
	struct pl_map		*map;
	struct pl_map_init_attr	 mia;
	int			 rc;

	D_RWLOCK_WRLOCK(&pl_rwlock);
	if (!pl_htable.ht_ops) {
		rc = pl_htable_init();
		if (rc)
			D_GOTO(out, rc);

		link = NULL;
	} else {
		/* already created */
		link = d_hash_rec_find(&pl_htable, uuid, sizeof(uuid_t));
	}

	if (!link) {
		pl_map_attr_init(pool_map, PL_TYPE_RING, &mia);
		rc = pl_map_create_inited(pool_map, &mia, &map);
		if (rc != 0)
			D_GOTO(out, rc);
	} else {
		struct pl_map	*tmp;

		tmp = container_of(link, struct pl_map, pl_link);
		if (pl_map_version(tmp) >= pool_map_get_version(pool_map)) {
			d_hash_rec_decref(&pl_htable, link);
			if (connect)
				tmp->pl_connects++;
			D_GOTO(out, rc = 0);
		}

		pl_map_attr_init(pool_map, PL_TYPE_RING, &mia);
		rc = pl_map_create_inited(pool_map, &mia, &map);
		if (rc != 0) {
			d_hash_rec_decref(&pl_htable, link);
			D_GOTO(out, rc);
		}

		/* transfer the pool connection count */
		map->pl_connects = tmp->pl_connects;
		/* evict the old placement map for this pool */
		d_hash_rec_delete_at(&pl_htable, link);
		d_hash_rec_decref(&pl_htable, link);
	}

	if (connect)
		map->pl_connects++;

	/* insert the new placement map into hash table */
	uuid_copy(map->pl_uuid, uuid);
	rc = d_hash_rec_insert(&pl_htable, uuid, sizeof(uuid_t),
			       &map->pl_link, true);
	D_ASSERT(rc == 0);
	pl_map_decref(map); /* hash table has held the refcount */
 out:
	D_RWLOCK_UNLOCK(&pl_rwlock);
	return rc;
}

/**
 * Drop the pool connection count of pl_map identified by \a uuid, evit it
 * from the placement map cache when connection count drops to zero.
 */
void
pl_map_disconnect(uuid_t uuid)
{
	d_list_t	*link;

	D_RWLOCK_WRLOCK(&pl_rwlock);
	link = d_hash_rec_find(&pl_htable, uuid, sizeof(uuid_t));
	if (link) {
		struct pl_map	*map;

		map = container_of(link, struct pl_map, pl_link);
		D_ASSERT(map->pl_connects > 0);
		map->pl_connects--;
		if (map->pl_connects == 0) {
			d_hash_rec_delete_at(&pl_htable, link);
			d_hash_rec_decref(&pl_htable, link);
		}
	}
	D_RWLOCK_UNLOCK(&pl_rwlock);
}

/**
 * Find the placement map of the pool identified by \a uuid.
 */
struct pl_map *
pl_map_find(uuid_t uuid, daos_obj_id_t oid)
{
	d_list_t	*link;

	D_RWLOCK_RDLOCK(&pl_rwlock);
	link = d_hash_rec_find(&pl_htable, uuid, sizeof(uuid_t));
	D_RWLOCK_UNLOCK(&pl_rwlock);

	return link ? pl_link2map(link) : NULL;
}
void
pl_map_addref(struct pl_map *map)
{
	d_hash_rec_addref(&pl_htable, &map->pl_link);
}

void
pl_map_decref(struct pl_map *map)
{
	d_hash_rec_decref(&pl_htable, &map->pl_link);
}

uint32_t
pl_map_version(struct pl_map *map)
{
	return map->pl_poolmap ? pool_map_get_version(map->pl_poolmap) : 0;
}
