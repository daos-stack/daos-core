/*
 * (C) Copyright 2019 Intel Corporation.
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
 * provided in Contract No. 8F-30005.
 * Any reproduction of computer software, computer software documentation, or
 * portions thereof marked with this legend must also reproduce the markings.
 */

#include <daos_types.h>
#include <daos_api.h>
#include <gurt/common.h>

#define DAOS_ACL_VERSION	1

/*
 * Comparison function for qsort. Compares by principal type.
 * Enum daos_acl_principal_type is in the expected order of type priority.
 */
static int
compare_aces(const void *p1, const void *p2)
{
	/* the inputs are in fact ptrs to ptrs */
	struct daos_ace *ace1 = *((struct daos_ace **)p1);
	struct daos_ace *ace2 = *((struct daos_ace **)p2);

	return (int)ace1->dae_principal_type - (int)ace2->dae_principal_type;
}

static void
sort_aces_by_principal_type(struct daos_ace *aces[], uint16_t num_aces)
{
	qsort(aces, num_aces, sizeof(struct daos_ace *), compare_aces);
}

/*
 * Flattens the array of ACE pointers into a single data blob in buffer.
 * Assumes caller has allocated the buffer large enough to hold the flattened
 * list.
 */
static void
flatten_aces(uint8_t *buffer, struct daos_ace *aces[], uint16_t num_aces)
{
	int	i;
	uint8_t	*current_ace;

	current_ace = buffer;
	for (i = 0; i < num_aces; i++) {
		int ace_size = daos_ace_get_size(aces[i]);

		memcpy(current_ace, aces[i], ace_size);
		current_ace += ace_size;
	}
}

/*
 * Calculates the expected length of the flattened ACE data blob.
 *
 * Returns -DER_INVAL if one of the ACEs is NULL.
 */
static int
get_flattened_ace_size(struct daos_ace *aces[], uint16_t num_aces)
{
	int	i;
	int	total_size = 0;

	for (i = 0; i < num_aces; i++) {
		int len = daos_ace_get_size(aces[i]);

		if (len < 0) {
			return len;
		}

		total_size += len;
	}

	return total_size;
}

struct daos_acl *
daos_acl_alloc(struct daos_ace *aces[], uint16_t num_aces)
{
	struct daos_acl	*acl;
	int		ace_len;

	ace_len = get_flattened_ace_size(aces, num_aces);
	if (ace_len < 0) {
		/* Bad ACE list */
		return NULL;
	}

	sort_aces_by_principal_type(aces, num_aces);

	D_ALLOC(acl, sizeof(struct daos_acl) + ace_len);
	if (acl == NULL) {
		/* Couldn't allocate */
		return NULL;
	}

	acl->dal_ver = DAOS_ACL_VERSION;
	acl->dal_len = ace_len;

	flatten_aces(acl->dal_ace, aces, num_aces);

	return acl;
}

void
daos_acl_free(struct daos_acl *acl)
{
	/* The ACL is one contiguous data blob - nothing special to do */
	D_FREE(acl);
}

static bool
principal_name_matches_ace(struct daos_ace *ace, const char *principal)
{
	if (principal == NULL) {
		/* Nothing to compare */
		return true;
	}

	return (strncmp(principal, ace->dae_principal, ace->dae_principal_len)
			== 0);
}

static bool
ace_matches_principal(struct daos_ace *ace,
		enum daos_acl_principal_type type, const char *principal,
		size_t principal_len)
{
	return	(ace->dae_principal_type == type) &&
		(ace->dae_principal_len == D_ALIGNUP(principal_len, 8)) &&
		principal_name_matches_ace(ace, principal);
}

static bool
principals_match(struct daos_ace *ace1, struct daos_ace *ace2)
{
	return ace_matches_principal(ace1, ace2->dae_principal_type,
			ace2->dae_principal, ace2->dae_principal_len);
}

/*
 * Write the ACE to the memory address pointed to by the pen.
 * Returns the new pen position.
 */
static uint8_t *
write_ace(struct daos_ace *ace, uint8_t *pen)
{
	size_t	len;

	len = daos_ace_get_size(ace);
	memcpy(pen, (uint8_t *)ace, len);

	return pen + len;
}

static void
copy_acl_with_new_ace_inserted(struct daos_acl *acl, struct daos_acl *new_acl,
		struct daos_ace *new_ace)
{
	struct daos_ace	*current;
	uint8_t		*pen;
	bool		new_written = false;

	current = daos_acl_get_first_ace(acl);
	pen = new_acl->dal_ace;
	while (current != NULL) {
		if (!new_written && current->dae_principal_type >
				new_ace->dae_principal_type)  {
			new_written = true;
			pen = write_ace(new_ace, pen);
		} else if (!new_written && principals_match(current,
					new_ace)) {
			new_written = true;
			pen = write_ace(new_ace, pen);

			/* new ACE already in old list - shorten the length */
			new_acl->dal_len = acl->dal_len;

			current = daos_acl_get_next_ace(acl, current);
		} else {
			pen = write_ace(current, pen);
			current = daos_acl_get_next_ace(acl, current);
		}
	}

	/* didn't insert it - tack it on the end */
	if (!new_written) {
		write_ace(new_ace, pen);
	}
}

int
daos_acl_add_ace_realloc(struct daos_acl *acl, struct daos_ace *new_ace,
		struct daos_acl **new_acl)
{
	int		new_ace_len;
	int		new_len;

	if (acl == NULL || new_acl == NULL) {
		return -DER_INVAL;
	}

	new_ace_len = daos_ace_get_size(new_ace);
	if (new_ace_len < 0) {
		/* ACE was invalid */
		return -DER_INVAL;
	}

	new_len = acl->dal_len + new_ace_len;

	D_ALLOC(*new_acl, sizeof(struct daos_acl) + new_len);
	if (*new_acl == NULL) {
		return -DER_NOMEM;
	}

	(*new_acl)->dal_ver = acl->dal_ver;
	(*new_acl)->dal_len = new_len;

	copy_acl_with_new_ace_inserted(acl, *new_acl, new_ace);

	return DER_SUCCESS;
}

static bool
type_is_valid(enum daos_acl_principal_type type)
{
	bool result = false;

	switch (type) {
	case DAOS_ACL_USER:
	case DAOS_ACL_GROUP:
	case DAOS_ACL_OWNER:
	case DAOS_ACL_OWNER_GROUP:
	case DAOS_ACL_EVERYONE:
		result = true;
		break;
	}

	return result;
}

static bool
type_needs_name(enum daos_acl_principal_type type)
{
	/*
	 * The only ACE types that require a name are User and Group. All others
	 * are "special" ACEs that apply to an abstract category.
	 */
	if (type == DAOS_ACL_USER || type == DAOS_ACL_GROUP) {
		return true;
	}

	return false;
}

static bool
principal_meets_type_requirements(enum daos_acl_principal_type type,
		const char *principal_name, size_t principal_name_len)
{
	return	(!type_needs_name(type) ||
		(principal_name != NULL && principal_name_len != 0));
}

int
daos_acl_remove_ace_realloc(struct daos_acl *acl,
		enum daos_acl_principal_type type, const char *principal_name,
		size_t principal_name_len, struct daos_acl **new_acl)
{
	struct daos_ace	*current;
	uint8_t		*pen;

	if (acl == NULL || new_acl == NULL || !type_is_valid(type) ||
		!principal_meets_type_requirements(type, principal_name,
				principal_name_len)) {
		return -DER_INVAL;
	}

	if (daos_acl_get_ace_for_principal(acl, type, principal_name) == NULL) {
		/* requested principal not in the list */
		return -DER_NONEXIST;
	}

	D_ALLOC(*new_acl, sizeof(struct daos_acl) + acl->dal_len);
	if (*new_acl == NULL) {
		return -DER_NOMEM;
	}

	(*new_acl)->dal_ver = acl->dal_ver;

	pen = (*new_acl)->dal_ace;
	current = daos_acl_get_first_ace(acl);
	while (current != NULL) {
		if (!ace_matches_principal(current, type, principal_name,
				principal_name_len)) {
			pen = write_ace(current, pen);
			(*new_acl)->dal_len += daos_ace_get_size(current);
		}

		current = daos_acl_get_next_ace(acl, current);
	}

	return DER_SUCCESS;
}

struct daos_ace *
daos_acl_get_first_ace(struct daos_acl *acl)
{
	if (acl == NULL || acl->dal_len == 0) {
		return NULL;
	}

	return (struct daos_ace *)acl->dal_ace;
}

static bool
is_in_ace_list(uint8_t *addr, struct daos_acl *acl)
{
	uint8_t *start_addr = acl->dal_ace;
	uint8_t *end_addr = acl->dal_ace + acl->dal_len;

	return addr >= start_addr && addr < end_addr;
}

struct daos_ace *
daos_acl_get_next_ace(struct daos_acl *acl, struct daos_ace *current_ace)
{
	size_t offset;

	if (acl == NULL || current_ace == NULL) {
		return NULL;
	}

	/* already at/beyond the end */
	if (!is_in_ace_list((uint8_t *)current_ace, acl)) {
		return NULL;
	}

	/* there is no next item */
	offset = sizeof(struct daos_ace) + current_ace->dae_principal_len;
	if (!is_in_ace_list((uint8_t *)current_ace + offset, acl)) {
		return NULL;
	}

	return (struct daos_ace *)((uint8_t *)current_ace + offset);
}

struct daos_ace *
daos_acl_get_ace_for_principal(struct daos_acl *acl,
		enum daos_acl_principal_type type, const char *principal)
{
	struct daos_ace *result;

	if (acl == NULL || !type_is_valid(type) ||
		(type_needs_name(type) && principal == NULL)) {
		return NULL;
	}

	result = daos_acl_get_first_ace(acl);
	while (result != NULL) {
		if (result->dae_principal_type == type &&
			principal_name_matches_ace(result, principal)) {
			break;
		}

		result = daos_acl_get_next_ace(acl, result);
	}

	return result;
}



static bool
type_is_group(enum daos_acl_principal_type type)
{
	if (type == DAOS_ACL_GROUP || type == DAOS_ACL_OWNER_GROUP) {
		return true;
	}

	return false;
}

struct daos_ace *
daos_ace_alloc(enum daos_acl_principal_type type, const char *principal_name,
		size_t principal_name_len)
{
	struct daos_ace	*ace;
	size_t		principal_array_len = 0;

	if (!type_is_valid(type)) {
		return NULL;
	}

	if (type_needs_name(type)) {
		if (principal_name == NULL || principal_name_len == 0) {
			return NULL;
		}

		/* align to 64 bits */
		principal_array_len = D_ALIGNUP(principal_name_len, 8);
	}

	D_ALLOC(ace, sizeof(struct daos_ace) + principal_array_len);
	if (ace != NULL) {
		ace->dae_principal_type = type;
		ace->dae_principal_len = principal_array_len;
		strncpy(ace->dae_principal, principal_name,
				principal_array_len);

		if (type_is_group(type)) {
			ace->dae_access_flags |= DAOS_ACL_FLAG_GROUP;
		}
	}

	return ace;
}


void
daos_ace_free(struct daos_ace *ace)
{
	D_FREE(ace);
}

int
daos_ace_get_size(struct daos_ace *ace)
{
	if (ace == NULL) {
		return -DER_INVAL;
	}

	return sizeof(struct daos_ace) + ace->dae_principal_len;
}
