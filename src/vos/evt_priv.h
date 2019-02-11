/**
 * (C) Copyright 2017-2019 Intel Corporation.
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
#ifndef __EVT_PRIV_H__
#define __EVT_PRIV_H__

#include <daos_srv/evtree.h>

/**
 * Tree node types.
 * NB: a node can be both root and leaf.
 */
enum {
	EVT_NODE_LEAF		= (1 << 0),	/**< leaf node */
	EVT_NODE_ROOT		= (1 << 1),	/**< root node */
};

enum evt_iter_state {
	EVT_ITER_NONE		= 0,	/**< uninitialized iterator */
	EVT_ITER_INIT,			/**< initialized but not probed */
	EVT_ITER_READY,			/**< probed, ready to iterate */
	EVT_ITER_FINI,			/**< reach at the end of iteration */
};

struct evt_iterator {
	/* Epoch range for the iterator */
	struct evt_filter		it_filter;
	/** state of the iterator */
	unsigned int			it_state;
	/** options for iterator */
	unsigned int			it_options;
	/** true if forward iterator */
	bool				it_forward;
	/** index */
	int				it_index;
	/** For sorted iterators */
	struct evt_entry_array		it_entries;
};

#define EVT_TRACE_MAX                   32

struct evt_trace {
	/** the current node offset */
	uint64_t			tr_node;
	/** child position of the searching trace */
	unsigned int			tr_at;
	/** Indicates whether node has been added to tx */
	bool				tr_tx_added;
};

struct evt_context {
	/** mapped address of the tree root */
	struct evt_root			*tc_root;
	/** memory ID of the tree root */
	TMMID(struct evt_root)		 tc_root_mmid;
	/** magic number to identify invalid tree open handle */
	unsigned int			 tc_magic;
	/** refcount on the context */
	unsigned int			 tc_ref;
	/** cached tree order (reduce PMEM access) */
	uint16_t			 tc_order;
	/** cached tree depth (reduce PMEM access) */
	uint16_t			 tc_depth;
	/** cached number of bytes per entry */
	uint32_t			 tc_inob;
	/** cached tree feature bits (reduce PMEM access) */
	uint64_t			 tc_feats;
	/** memory instance (PMEM or DRAM) */
	struct umem_instance		 tc_umm;
	/** pmemobj pool uuid */
	uint64_t			 tc_pmempool_uuid;
	/**
	 * NVMe free space tracking information used for NVMe record
	 * alloc & free
	 */
	void				*tc_blks_info;
	/** embedded iterator */
	struct evt_iterator		 tc_iter;
	/** space to store tree search path */
	struct evt_trace		 tc_trace_scratch[EVT_TRACE_MAX];
	/** points to &tc_trace_scratch[EVT_TRACE_MAX - depth] */
	struct evt_trace		*tc_trace;
	/** customized operation table for different tree policies */
	struct evt_policy_ops		*tc_ops;
	/* The container open handle */
	daos_handle_t			 tc_coh;
};

#define EVT_NODE_NULL			TMMID_NULL(struct evt_node)
#define EVT_ROOT_NULL			TMMID_NULL(struct evt_root)

#define evt_umm(tcx)			(&(tcx)->tc_umm)
#define evt_tmmid2ptr(tcx, tmmid)	umem_id2ptr_typed(evt_umm(tcx), tmmid)
#define evt_has_tx(tcx)			umem_has_tx(evt_umm(tcx))

#define evt_off2mmid(tcx, offset)			\
(							\
{							\
	umem_id_t	__ummid;			\
							\
	__ummid.pool_uuid_lo = (tcx)->tc_pmempool_uuid;	\
	__ummid.off = (offset);				\
	__ummid;					\
}							\
)

#define evt_off2tmmid(tcx, offset, type)		\
	umem_id_u2t(evt_off2mmid(tcx, offset), type)

#define evt_off2ptr(tcx, offset)			\
	umem_id2ptr(evt_umm(tcx), evt_off2mmid(tcx, offset))

#define evt_off2nodemmid(tcx, offset)			\
	evt_off2tmmid(tcx, offset, struct evt_node)

#define evt_off2descmmid(tcx, offset)			\
	evt_off2tmmid(tcx, offset, struct evt_desc)

#define EVT_NODE_MAGIC 0xfefef00d
#define EVT_DESC_MAGIC 0xbeefdead

/** Convert an offset to a evtree node descriptor
 * \param[IN]	tcx	Tree context
 * \param[IN]	offset	The offset in the umem pool
 */
static inline struct evt_node *
evt_off2node(struct evt_context *tcx, uint64_t offset)
{
	struct evt_node *node;

	node = evt_off2ptr(tcx, offset);
	D_ASSERT(node->tn_magic == EVT_NODE_MAGIC);

	return node;
}

/** Convert an offset to a evtree data descriptor
 * \param[IN]	tcx	Tree context
 * \param[IN]	offset	The offset in the umem pool
 */
static inline struct evt_desc *
evt_off2desc(struct evt_context *tcx, uint64_t offset)
{
	struct evt_desc *desc;

	desc = evt_off2ptr(tcx, offset);
	D_ASSERT(desc->dc_magic == EVT_DESC_MAGIC);

	return desc;
}

/** Helper function for starting a PMDK transaction, if applicable */
static inline int
evt_tx_begin(struct evt_context *tcx)
{
	if (!evt_has_tx(tcx))
		return 0;

	return umem_tx_begin(evt_umm(tcx), NULL);
}

/** Helper function for ending a PMDK transaction, if applicable */
static inline int
evt_tx_end(struct evt_context *tcx, int rc)
{
	if (!evt_has_tx(tcx))
		return rc;

	if (rc != 0)
		return umem_tx_abort(evt_umm(tcx), rc);

	return umem_tx_commit(evt_umm(tcx));
}

/* By definition, all rectangles overlap in the epoch range because all
 * are from start to infinity.  However, for common queries, we often only want
 * rectangles intersect at a given epoch
 */
enum evt_find_opc {
	/** find all rectangles overlapped with the input rectangle */
	EVT_FIND_ALL,
	/** find the first rectangle overlapped with the input rectangle */
	EVT_FIND_FIRST,
	/**
	 * Returns -DER_NO_PERM if any overlapping rectangle is found in
	 * the same epoch with an identical sequence number.
	 */
	EVT_FIND_OVERWRITE,
	/** Find the exactly same extent. */
	EVT_FIND_SAME,
};

/** Clone an evtree context
 *
 * \param[IN]	tcx	The context to clone
 * \param[OUT]	tcx_pp	The new context
 */
int evt_tcx_clone(struct evt_context *tcx, struct evt_context **tcx_pp);

/** Remove the leaf node at the current trace
 *
 * \param[IN]	tcx	The context to use for delete
 * \param[IN]	remove	If true, payload is free'd internally
 *
 *  The trace is set as if evt_move_trace were called.
 *
 * Returns -DER_NONEXIST if it's the last item in the trace
 */
int evt_node_delete(struct evt_context *tcx, bool remove);

#define EVT_HDL_ALIVE	0xbabecafe
#define EVT_HDL_DEAD	0xdeadbeef

static inline void
evt_tcx_addref(struct evt_context *tcx)
{
	tcx->tc_ref++;
	if (tcx->tc_inob == 0)
		tcx->tc_inob = tcx->tc_root->tr_inob;
}

static inline void
evt_tcx_decref(struct evt_context *tcx)
{
	D_ASSERT(tcx->tc_ref > 0);
	tcx->tc_ref--;
	if (tcx->tc_ref == 0) {
		tcx->tc_magic = EVT_HDL_DEAD;
		/* Free any memory allocated by embedded iterator */
		evt_ent_array_fini(&tcx->tc_iter.it_entries);
		D_FREE(tcx);
	}
}

/** Return true if a rectangle doesn't intersect the filter
 *
 * \param[IN]	filter	The optional input filter
 * \param[IN]	rect	The rectangle to check
 * \param[IN]	leaf	Indicates if the rectangle is a leaf entry
 */
static inline bool
evt_filter_rect(const struct evt_filter *filter, const struct evt_rect *rect,
		bool leaf)
{
	if (filter == NULL)
		goto done;

	if (filter->fr_ex.ex_lo > rect->rc_ex.ex_hi ||
	    filter->fr_ex.ex_hi < rect->rc_ex.ex_lo ||
	    filter->fr_epr.epr_hi < rect->rc_epc)
		return true; /* Rectangle is outside of filter */

	/* In tree rectangle only includes lower bound.  For intermediate
	 * nodes, we can't filter based on lower bound.  For leaf nodes,
	 * we can because it represents a point in time.
	 */
	if (!leaf)
		goto done;

	if (filter->fr_epr.epr_lo > rect->rc_epc)
		return true; /* Rectangle is outside of filter */
done:
	return false;
}

/** Create an equivalent evt_rect from an evt_entry
 *
 * \param[OUT]	rect	The output rectangle
 * \param[IN]	ent	The input entry
 */
static inline void
evt_ent2rect(struct evt_rect *rect, const struct evt_entry *ent)
{
	rect->rc_ex = ent->en_sel_ext;
	rect->rc_epc = ent->en_epoch;
}

/** Sort entries in an entry array
 * \param[IN]		tcx		The evtree context
 * \param[IN, OUT]	ent_array	The entry array to sort
 * \param[IN]		flags		Visibility flags
 *					EVT_VISIBLE: Return visible records
 *					EVT_COVERED: Return covered records
 * Returns 0 if successful, error otherwise.   The resulting array will
 * be sorted by start offset, high epoch
 */
int evt_ent_array_sort(struct evt_context *tcx,
		       struct evt_entry_array *ent_array, int flags);
/** Scan the tree and select all rectangles that match
 * \param[IN]		tcx		The evtree context
 * \param[IN]		opc		The opcode for the scan
 * \param[IN]		intent		The operation intent
 *					EVT_FIND_FIRST: First record only
 *					EVT_FIND_SAME:  Same record only
 *					EVT_FIND_ALL:   All records
 * \param[IN]		filter		Filters for records
 * \param[IN]		rect		The specific rectangle to match
 * \param[IN,OUT]	ent_array	The initialized array to fill
 *
 * Returns 0 if successful, error otherwise.  The tree trace will point at last
 * scanned record.
 */
int evt_ent_array_fill(struct evt_context *tcx, enum evt_find_opc find_opc,
		       uint32_t intent, const struct evt_filter *filter,
		       const struct evt_rect *rect,
		       struct evt_entry_array *ent_array);

/** Compare two rectanglesConvert a context to a daos_handle
 * \param[IN]		rt1	The first rectangle
 * \param[IN]		rt2	The second rectangle
 *
 * Returns < 0 if rt1 < rt2
 * returns > 0 if rt1 > rt2
 * returns 0 if equal
 *
 * Order is lower high offset, higher epoch, lower high offset
 */
int evt_rect_cmp(const struct evt_rect *rt1, const struct evt_rect *rt2);

/** Convert a context to a daos_handle
 * \param[IN]	tcx	The evtree context
 *
 * Returns the converted handle
 */
daos_handle_t evt_tcx2hdl(struct evt_context *tcx);

/** Convert a handle to an evtree context
 * \param[IN]	handle	The daos handle
 *
 * Returns the converted handle
 */
struct evt_context *evt_hdl2tcx(daos_handle_t toh);

/** Move the trace forward.
 * \param[IN]	tcx	The evtree context
 * \param IN]	intent	The operation intent
 */
bool evt_move_trace(struct evt_context *tcx, uint32_t intent);

/** Get a pointer to the rectangle corresponding to an index in a tree node
 * \param[IN]	tcx	The evtree context
 * \param[IN]	nd_mmid	The tree node
 * \param[IN]	at	The index in the node
 *
 * Returns the rectangle at the index
 */
struct evt_rect *evt_node_rect_at(struct evt_context *tcx,
				  uint64_t nd_off, unsigned int at);

/** Fill an evt_entry from the record at an index in a tree node
 * \param[IN]	tcx		The evtree context
 * \param[IN]	nd_mmid		The tree node
 * \param[IN]	at		The index in the node
 * \param[IN]	rect_srch	The original rectangle used for the search
 * \param[OUT]	entry		The entry to fill
 *
 * The selected extent will be trimmed by the search rectangle used.
 */
void evt_entry_fill(struct evt_context *tcx, uint64_t nd_off,
		    unsigned int at, const struct evt_rect *rect_srch,
		    struct evt_entry *entry);
#endif /* __EVT_PRIV_H__ */
