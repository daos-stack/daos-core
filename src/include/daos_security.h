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

/**
 * \file
 *
 * DAOS API methods for security and access control
 */

#ifndef __DAOS_SECURITY_H__
#define __DAOS_SECURITY_H__

#if defined(__cplusplus)
extern "C" {
#endif

#include <stdint.h>
#include <stdbool.h>
#include <sys/types.h>

#define	DAOS_ACL_VERSION		1

/**
 * Header for the Access Control List, followed by the table of variable-length
 * Access Control Entries.
 * The entry list may be walked by inspecting the principal length and
 * calculating the entry's overall length from that.
 */
struct daos_acl {
	/** Version of the table format */
	uint16_t	dal_ver;
	/** reserved for 64-bit alignment */
	uint16_t	dal_reserv;
	/** length of entries list in bytes */
	uint32_t	dal_len;
	/** table of variable-length Access Control Entries (struct daos_ace) */
	uint8_t		dal_ace[];
};

/**
 * Type of principal for the Access Control Entry.
 * OWNER, OWNER_GROUP, and EVERYONE are special principals that do not need
 * a principal name string.
 */
enum daos_acl_principal_type {
	DAOS_ACL_OWNER,		/** Owner of the object */
	DAOS_ACL_USER,		/** Individual user */
	DAOS_ACL_OWNER_GROUP,	/** Owning group */
	DAOS_ACL_GROUP,		/** Group */
	DAOS_ACL_EVERYONE,	/** Anyone else */

	NUM_DAOS_ACL_TYPES	/** Must be last */
};

/**
 * Bits representing access types to set permissions for
 */
enum daos_acl_access_type {
	DAOS_ACL_ACCESS_ALLOW = (1U << 0),	/** allow access */
	DAOS_ACL_ACCESS_AUDIT = (1U << 1),	/** log the access for review */
	DAOS_ACL_ACCESS_ALARM = (1U << 2)	/** notify of the access */
};

/**
 * Bits representing access flags
 */
enum daos_acl_flags {
	/** This represents a group, not a user */
	DAOS_ACL_FLAG_GROUP		= (1U << 0),
	/** Containers should inherit access controls from this pool */
	DAOS_ACL_FLAG_POOL_INHERIT	= (1U << 1),
	/** Audit/alarm should occur on failed access */
	DAOS_ACL_FLAG_ACCESS_FAIL	= (1U << 2),
	/** Audit/alarm should occur on successful access */
	DAOS_ACL_FLAG_ACCESS_SUCCESS	= (1U << 3)
};

/**
 * Bits representing the specific permissions that may be set
 */
enum daos_acl_perm {
	DAOS_ACL_PERM_READ	= (1U << 0),
	DAOS_ACL_PERM_WRITE	= (1U << 1)
};

/**
 * Access Control Entry for a given principal.
 * Each principal has at most one ACE that lists all their permissions in a
 * given Access Control List.
 */
struct daos_ace {
	/** Bitmap of daos_acl_access_type */
	uint8_t		dae_access_types;
	/** daos_acl_principal_type */
	uint8_t		dae_principal_type;
	/** Length of the principal string */
	uint16_t	dae_principal_len;
	/** Bitmap of daos_acl_flags */
	uint16_t	dae_access_flags;
	/** Reserved for 64-bit alignment */
	uint16_t	dae_reserv;
	/** Bitmap of daos_acl_perm for the ALLOW access */
	uint64_t	dae_allow_perms;
	/** Bitmap of daos_acl_perm for AUDIT access */
	uint64_t	dae_audit_perms;
	/** Bitmap of daos_acl_perm for ALARM access */
	uint64_t	dae_alarm_perms;
	/**
	 * Null-terminated string representing the principal name for specific
	 * user/group.
	 * Actual bytes allocated MUST be rounded up for 64-bit alignment.
	 * Empty for special principals OWNER, OWNER_GROUP, and EVERYONE.
	 */
	char		dae_principal[];
};

/**
 * Allocate an DAOS Access Control List.
 *
 * \param[in]	aces		Array of pointers to ACEs to be put in the ACL.
 * \param[in]	num_aces	Number of ACEs in array
 *
 * \return	allocated daos_acl pointer, NULL if failed
 */
struct daos_acl *
daos_acl_create(struct daos_ace *aces[], uint16_t num_aces);

/**
 * Allocate a new copy of a DAOS Access Control List.
 *
 * \param[in]	acl	ACL structure to be copied
 *
 * \return	Newly allocated copy of the ACL, or NULL if the ACL can't be
 *		allocated
 */
struct daos_acl *
daos_acl_copy(struct daos_acl *acl);

/**
 * Free a DAOS Access Control List.
 *
 * \param[in]	acl	ACL pointer to be freed
 */
void
daos_acl_free(struct daos_acl *acl);

/**
 * Get the next Access Control Entry in the Access Control List, for iterating
 * over the list.
 *
 * \param[in]	acl		ACL to traverse
 * \param[in]	current_ace	Current ACE, to determine the next one, or NULL
 *				for the first ACE
 *
 * \return	Pointer to the next ACE in the ACL, or NULL if at the end
 */
struct daos_ace *
daos_acl_get_next_ace(struct daos_acl *acl, struct daos_ace *current_ace);

/**
 * Search the Access Control List for an Access Control Entry for a specific
 * principal.
 *
 * \param[in]	acl		ACL to search
 * \param[in]	type		Principal type to search for
 * \param[in]	principal	Principal name, if type is USER or GROUP. NULL
 *				otherwise.
 * \param[out]	ace		Pointer to matching ACE within ACL (not a copy)
 *
 * \return	0		Success
 *		-DER_INVAL	Invalid input
 *		-DER_NONEXIST	Matching ACE not found
 */
int
daos_acl_get_ace_for_principal(struct daos_acl *acl,
			       enum daos_acl_principal_type type,
			       const char *principal, struct daos_ace **ace);

/**
 * Insert an Access Control Entry in the appropriate location in the ACE
 * list. The expected order is: Owner, Users, Assigned Group, Groups, Everyone.
 *
 * The ACL structure may be reallocated to make room for the new ACE. If so the
 * old structure will be freed.
 *
 * If the new ACE is an update of an existing entry, it will replace the old
 * entry.
 *
 * \param[in]	acl	ACL to modify
 * \param[in]	new_ace	ACE to be added
 *
 * \return	0		Success
 *		-DER_INVAL	Invalid input
 *		-DER_NOMEM	Failed to allocate required memory
 */
int
daos_acl_add_ace(struct daos_acl **acl, struct daos_ace *new_ace);

/**
 * Remove an Access Control Entry from the list.
 *
 * When the entry is removed, the ACL is reallocated, and the old structure
 * is freed.
 *
 * \param[in]	acl			Original ACL
 * \param[in]	type			Principal type of the ACE to remove
 * \param[in]	principal_name		Principal name of the ACE to remove
 *					(NULL if type isn't user/group)
 * \param[out]	new_acl			Reallocated copy of the ACL with the
 *					ACE removed
 *
 * \return	0		Success
 *		-DER_INVAL	Invalid input
 *		-DER_NOMEM	Failed to allocate required memory
 *		-DER_NONEXIST	Requested ACE was not in the ACL
 */
int
daos_acl_remove_ace(struct daos_acl **acl,
		    enum daos_acl_principal_type type,
		    const char *principal_name);

/**
 * Print the Access Control List to stdout in a human-readable format.
 *
 * \param	acl	Access Control List to print
 */
void
daos_acl_dump(struct daos_acl *acl);

/**
 * Parse and sanity check the entire Access Control List for valid values and
 * internal consistency.
 *
 * \param	acl	Access Control List to sanity check
 *
 * \return	True if the ACL is valid, false otherwise
 */
bool
daos_acl_is_valid(struct daos_acl *acl);

/**
 * Allocate a new Access Control Entry with an appropriately aligned principal
 * name, if applicable.
 *
 * Only User and Group types use principal name.
 *
 * \param[in]	type			Type of principal for the ACE
 * \param[in]	principal_name		Principal name will be added to the end
 *					of the structure. For types that don't
 *					use it, it is ignored. OK to pass NULL.
 *
 * \return	New ACE structure with an appropriately packed principal name,
 *			length, and type set.
 */
struct daos_ace *
daos_ace_create(enum daos_acl_principal_type type, const char *principal_name);

/**
 * Free an Access Control Entry allocated by daos_ace_alloc().
 *
 * \param[in]	ace	ACE to be freed
 */
void
daos_ace_free(struct daos_ace *ace);

/**
 * Get the length in bytes of an Access Control Entry.
 * The entries have variable length.
 *
 * \param[in]	ace	ACE to get the size of
 *
 * \return	Success		Size of ACE in bytes
 *		-DER_INVAL	Invalid input
 */
ssize_t
daos_ace_get_size(struct daos_ace *ace);

/**
 * Print the Access Control Entry to stdout in a human-readable format.
 *
 * \param	ace	Access Control Entry to print
 * \param	tabs	Number of tabs to indent at top level
 */
void
daos_ace_dump(struct daos_ace *ace, uint tabs);

/**
 * Sanity check the Access Control Entry structure for valid values and internal
 * consistency.
 *
 * \param	ace	Access Control Entry to be checked
 *
 * \return	True if the ACE is valid, false otherwise
 */
bool
daos_ace_is_valid(struct daos_ace *ace);

#if defined(__cplusplus)
}
#endif
#endif /* __DAOS_SECURITY_H__ */
