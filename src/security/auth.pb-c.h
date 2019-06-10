/* Generated by the protocol buffer compiler.  DO NOT EDIT! */
/* Generated from: auth.proto */

#ifndef PROTOBUF_C_auth_2eproto__INCLUDED
#define PROTOBUF_C_auth_2eproto__INCLUDED

#include <protobuf-c/protobuf-c.h>

PROTOBUF_C__BEGIN_DECLS

#if PROTOBUF_C_VERSION_NUMBER < 1003000
# error This file was generated by a newer version of protoc-c which is incompatible with your libprotobuf-c headers. Please update your headers.
#elif 1003001 < PROTOBUF_C_MIN_COMPILER_VERSION
# error This file was generated by an older version of protoc-c which is incompatible with your libprotobuf-c headers. Please regenerate this file with a newer version of protoc-c.
#endif


typedef struct _Auth__Token Auth__Token;
typedef struct _Auth__Sys Auth__Sys;
typedef struct _Auth__SysVerifier Auth__SysVerifier;
typedef struct _Auth__Credential Auth__Credential;


/* --- enums --- */

/*
 * Authentication token includes a packed structure of the specified flavor
 */
typedef enum _Auth__Flavor {
  AUTH__FLAVOR__AUTH_NONE = 0,
  AUTH__FLAVOR__AUTH_SYS = 1
    PROTOBUF_C__FORCE_ENUM_TO_BE_INT_SIZE(AUTH__FLAVOR)
} Auth__Flavor;

/* --- messages --- */

struct  _Auth__Token
{
  ProtobufCMessage base;
  Auth__Flavor flavor;
  ProtobufCBinaryData data;
};
#define AUTH__TOKEN__INIT \
 { PROTOBUF_C_MESSAGE_INIT (&auth__token__descriptor) \
    , AUTH__FLAVOR__AUTH_NONE, {0,NULL} }


/*
 * Token structure for AUTH_SYS flavor
 */
struct  _Auth__Sys
{
  ProtobufCMessage base;
  uint64_t stamp;
  char *machinename;
  char *user;
  char *group;
  size_t n_groups;
  char **groups;
  /*
   * Additional field for MAC label
   */
  char *secctx;
};
#define AUTH__SYS__INIT \
 { PROTOBUF_C_MESSAGE_INIT (&auth__sys__descriptor) \
    , 0, (char *)protobuf_c_empty_string, (char *)protobuf_c_empty_string, (char *)protobuf_c_empty_string, 0,NULL, (char *)protobuf_c_empty_string }


struct  _Auth__SysVerifier
{
  ProtobufCMessage base;
  ProtobufCBinaryData signature;
};
#define AUTH__SYS_VERIFIER__INIT \
 { PROTOBUF_C_MESSAGE_INIT (&auth__sys_verifier__descriptor) \
    , {0,NULL} }


/*
 * SecurityCredential includes the auth token and a verifier that can be used by
 * the server to verify the integrity of the token.
 * Token and verifier are expected to have the same flavor type.
 */
struct  _Auth__Credential
{
  ProtobufCMessage base;
  Auth__Token *token;
  Auth__Token *verifier;
};
#define AUTH__CREDENTIAL__INIT \
 { PROTOBUF_C_MESSAGE_INIT (&auth__credential__descriptor) \
    , NULL, NULL }


/* Auth__Token methods */
void   auth__token__init
                     (Auth__Token         *message);
size_t auth__token__get_packed_size
                     (const Auth__Token   *message);
size_t auth__token__pack
                     (const Auth__Token   *message,
                      uint8_t             *out);
size_t auth__token__pack_to_buffer
                     (const Auth__Token   *message,
                      ProtobufCBuffer     *buffer);
Auth__Token *
       auth__token__unpack
                     (ProtobufCAllocator  *allocator,
                      size_t               len,
                      const uint8_t       *data);
void   auth__token__free_unpacked
                     (Auth__Token *message,
                      ProtobufCAllocator *allocator);
/* Auth__Sys methods */
void   auth__sys__init
                     (Auth__Sys         *message);
size_t auth__sys__get_packed_size
                     (const Auth__Sys   *message);
size_t auth__sys__pack
                     (const Auth__Sys   *message,
                      uint8_t             *out);
size_t auth__sys__pack_to_buffer
                     (const Auth__Sys   *message,
                      ProtobufCBuffer     *buffer);
Auth__Sys *
       auth__sys__unpack
                     (ProtobufCAllocator  *allocator,
                      size_t               len,
                      const uint8_t       *data);
void   auth__sys__free_unpacked
                     (Auth__Sys *message,
                      ProtobufCAllocator *allocator);
/* Auth__SysVerifier methods */
void   auth__sys_verifier__init
                     (Auth__SysVerifier         *message);
size_t auth__sys_verifier__get_packed_size
                     (const Auth__SysVerifier   *message);
size_t auth__sys_verifier__pack
                     (const Auth__SysVerifier   *message,
                      uint8_t             *out);
size_t auth__sys_verifier__pack_to_buffer
                     (const Auth__SysVerifier   *message,
                      ProtobufCBuffer     *buffer);
Auth__SysVerifier *
       auth__sys_verifier__unpack
                     (ProtobufCAllocator  *allocator,
                      size_t               len,
                      const uint8_t       *data);
void   auth__sys_verifier__free_unpacked
                     (Auth__SysVerifier *message,
                      ProtobufCAllocator *allocator);
/* Auth__Credential methods */
void   auth__credential__init
                     (Auth__Credential         *message);
size_t auth__credential__get_packed_size
                     (const Auth__Credential   *message);
size_t auth__credential__pack
                     (const Auth__Credential   *message,
                      uint8_t             *out);
size_t auth__credential__pack_to_buffer
                     (const Auth__Credential   *message,
                      ProtobufCBuffer     *buffer);
Auth__Credential *
       auth__credential__unpack
                     (ProtobufCAllocator  *allocator,
                      size_t               len,
                      const uint8_t       *data);
void   auth__credential__free_unpacked
                     (Auth__Credential *message,
                      ProtobufCAllocator *allocator);
/* --- per-message closures --- */

typedef void (*Auth__Token_Closure)
                 (const Auth__Token *message,
                  void *closure_data);
typedef void (*Auth__Sys_Closure)
                 (const Auth__Sys *message,
                  void *closure_data);
typedef void (*Auth__SysVerifier_Closure)
                 (const Auth__SysVerifier *message,
                  void *closure_data);
typedef void (*Auth__Credential_Closure)
                 (const Auth__Credential *message,
                  void *closure_data);

/* --- services --- */


/* --- descriptors --- */

extern const ProtobufCEnumDescriptor    auth__flavor__descriptor;
extern const ProtobufCMessageDescriptor auth__token__descriptor;
extern const ProtobufCMessageDescriptor auth__sys__descriptor;
extern const ProtobufCMessageDescriptor auth__sys_verifier__descriptor;
extern const ProtobufCMessageDescriptor auth__credential__descriptor;

PROTOBUF_C__END_DECLS


#endif  /* PROTOBUF_C_auth_2eproto__INCLUDED */
