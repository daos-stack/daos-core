/* Generated by the protocol buffer compiler.  DO NOT EDIT! */
/* Generated from: auth.proto */

/* Do not generate deprecated warnings for self */
#ifndef PROTOBUF_C__NO_DEPRECATED
#define PROTOBUF_C__NO_DEPRECATED
#endif

#include "auth.pb-c.h"
void   auth__token__init
                     (Auth__Token         *message)
{
  static const Auth__Token init_value = AUTH__TOKEN__INIT;
  *message = init_value;
}
size_t auth__token__get_packed_size
                     (const Auth__Token *message)
{
  assert(message->base.descriptor == &auth__token__descriptor);
  return protobuf_c_message_get_packed_size ((const ProtobufCMessage*)(message));
}
size_t auth__token__pack
                     (const Auth__Token *message,
                      uint8_t       *out)
{
  assert(message->base.descriptor == &auth__token__descriptor);
  return protobuf_c_message_pack ((const ProtobufCMessage*)message, out);
}
size_t auth__token__pack_to_buffer
                     (const Auth__Token *message,
                      ProtobufCBuffer *buffer)
{
  assert(message->base.descriptor == &auth__token__descriptor);
  return protobuf_c_message_pack_to_buffer ((const ProtobufCMessage*)message, buffer);
}
Auth__Token *
       auth__token__unpack
                     (ProtobufCAllocator  *allocator,
                      size_t               len,
                      const uint8_t       *data)
{
  return (Auth__Token *)
     protobuf_c_message_unpack (&auth__token__descriptor,
                                allocator, len, data);
}
void   auth__token__free_unpacked
                     (Auth__Token *message,
                      ProtobufCAllocator *allocator)
{
  if(!message)
    return;
  assert(message->base.descriptor == &auth__token__descriptor);
  protobuf_c_message_free_unpacked ((ProtobufCMessage*)message, allocator);
}
void   auth__sys__init
                     (Auth__Sys         *message)
{
  static const Auth__Sys init_value = AUTH__SYS__INIT;
  *message = init_value;
}
size_t auth__sys__get_packed_size
                     (const Auth__Sys *message)
{
  assert(message->base.descriptor == &auth__sys__descriptor);
  return protobuf_c_message_get_packed_size ((const ProtobufCMessage*)(message));
}
size_t auth__sys__pack
                     (const Auth__Sys *message,
                      uint8_t       *out)
{
  assert(message->base.descriptor == &auth__sys__descriptor);
  return protobuf_c_message_pack ((const ProtobufCMessage*)message, out);
}
size_t auth__sys__pack_to_buffer
                     (const Auth__Sys *message,
                      ProtobufCBuffer *buffer)
{
  assert(message->base.descriptor == &auth__sys__descriptor);
  return protobuf_c_message_pack_to_buffer ((const ProtobufCMessage*)message, buffer);
}
Auth__Sys *
       auth__sys__unpack
                     (ProtobufCAllocator  *allocator,
                      size_t               len,
                      const uint8_t       *data)
{
  return (Auth__Sys *)
     protobuf_c_message_unpack (&auth__sys__descriptor,
                                allocator, len, data);
}
void   auth__sys__free_unpacked
                     (Auth__Sys *message,
                      ProtobufCAllocator *allocator)
{
  if(!message)
    return;
  assert(message->base.descriptor == &auth__sys__descriptor);
  protobuf_c_message_free_unpacked ((ProtobufCMessage*)message, allocator);
}
void   auth__sys_verifier__init
                     (Auth__SysVerifier         *message)
{
  static const Auth__SysVerifier init_value = AUTH__SYS_VERIFIER__INIT;
  *message = init_value;
}
size_t auth__sys_verifier__get_packed_size
                     (const Auth__SysVerifier *message)
{
  assert(message->base.descriptor == &auth__sys_verifier__descriptor);
  return protobuf_c_message_get_packed_size ((const ProtobufCMessage*)(message));
}
size_t auth__sys_verifier__pack
                     (const Auth__SysVerifier *message,
                      uint8_t       *out)
{
  assert(message->base.descriptor == &auth__sys_verifier__descriptor);
  return protobuf_c_message_pack ((const ProtobufCMessage*)message, out);
}
size_t auth__sys_verifier__pack_to_buffer
                     (const Auth__SysVerifier *message,
                      ProtobufCBuffer *buffer)
{
  assert(message->base.descriptor == &auth__sys_verifier__descriptor);
  return protobuf_c_message_pack_to_buffer ((const ProtobufCMessage*)message, buffer);
}
Auth__SysVerifier *
       auth__sys_verifier__unpack
                     (ProtobufCAllocator  *allocator,
                      size_t               len,
                      const uint8_t       *data)
{
  return (Auth__SysVerifier *)
     protobuf_c_message_unpack (&auth__sys_verifier__descriptor,
                                allocator, len, data);
}
void   auth__sys_verifier__free_unpacked
                     (Auth__SysVerifier *message,
                      ProtobufCAllocator *allocator)
{
  if(!message)
    return;
  assert(message->base.descriptor == &auth__sys_verifier__descriptor);
  protobuf_c_message_free_unpacked ((ProtobufCMessage*)message, allocator);
}
void   auth__credential__init
                     (Auth__Credential         *message)
{
  static const Auth__Credential init_value = AUTH__CREDENTIAL__INIT;
  *message = init_value;
}
size_t auth__credential__get_packed_size
                     (const Auth__Credential *message)
{
  assert(message->base.descriptor == &auth__credential__descriptor);
  return protobuf_c_message_get_packed_size ((const ProtobufCMessage*)(message));
}
size_t auth__credential__pack
                     (const Auth__Credential *message,
                      uint8_t       *out)
{
  assert(message->base.descriptor == &auth__credential__descriptor);
  return protobuf_c_message_pack ((const ProtobufCMessage*)message, out);
}
size_t auth__credential__pack_to_buffer
                     (const Auth__Credential *message,
                      ProtobufCBuffer *buffer)
{
  assert(message->base.descriptor == &auth__credential__descriptor);
  return protobuf_c_message_pack_to_buffer ((const ProtobufCMessage*)message, buffer);
}
Auth__Credential *
       auth__credential__unpack
                     (ProtobufCAllocator  *allocator,
                      size_t               len,
                      const uint8_t       *data)
{
  return (Auth__Credential *)
     protobuf_c_message_unpack (&auth__credential__descriptor,
                                allocator, len, data);
}
void   auth__credential__free_unpacked
                     (Auth__Credential *message,
                      ProtobufCAllocator *allocator)
{
  if(!message)
    return;
  assert(message->base.descriptor == &auth__credential__descriptor);
  protobuf_c_message_free_unpacked ((ProtobufCMessage*)message, allocator);
}
static const ProtobufCFieldDescriptor auth__token__field_descriptors[2] =
{
  {
    "flavor",
    1,
    PROTOBUF_C_LABEL_NONE,
    PROTOBUF_C_TYPE_ENUM,
    0,   /* quantifier_offset */
    offsetof(Auth__Token, flavor),
    &auth__flavor__descriptor,
    NULL,
    0,             /* flags */
    0,NULL,NULL    /* reserved1,reserved2, etc */
  },
  {
    "data",
    2,
    PROTOBUF_C_LABEL_NONE,
    PROTOBUF_C_TYPE_BYTES,
    0,   /* quantifier_offset */
    offsetof(Auth__Token, data),
    NULL,
    NULL,
    0,             /* flags */
    0,NULL,NULL    /* reserved1,reserved2, etc */
  },
};
static const unsigned auth__token__field_indices_by_name[] = {
  1,   /* field[1] = data */
  0,   /* field[0] = flavor */
};
static const ProtobufCIntRange auth__token__number_ranges[1 + 1] =
{
  { 1, 0 },
  { 0, 2 }
};
const ProtobufCMessageDescriptor auth__token__descriptor =
{
  PROTOBUF_C__MESSAGE_DESCRIPTOR_MAGIC,
  "auth.Token",
  "Token",
  "Auth__Token",
  "auth",
  sizeof(Auth__Token),
  2,
  auth__token__field_descriptors,
  auth__token__field_indices_by_name,
  1,  auth__token__number_ranges,
  (ProtobufCMessageInit) auth__token__init,
  NULL,NULL,NULL    /* reserved[123] */
};
static const ProtobufCFieldDescriptor auth__sys__field_descriptors[6] =
{
  {
    "stamp",
    1,
    PROTOBUF_C_LABEL_NONE,
    PROTOBUF_C_TYPE_UINT64,
    0,   /* quantifier_offset */
    offsetof(Auth__Sys, stamp),
    NULL,
    NULL,
    0,             /* flags */
    0,NULL,NULL    /* reserved1,reserved2, etc */
  },
  {
    "machinename",
    2,
    PROTOBUF_C_LABEL_NONE,
    PROTOBUF_C_TYPE_STRING,
    0,   /* quantifier_offset */
    offsetof(Auth__Sys, machinename),
    NULL,
    &protobuf_c_empty_string,
    0,             /* flags */
    0,NULL,NULL    /* reserved1,reserved2, etc */
  },
  {
    "user",
    3,
    PROTOBUF_C_LABEL_NONE,
    PROTOBUF_C_TYPE_STRING,
    0,   /* quantifier_offset */
    offsetof(Auth__Sys, user),
    NULL,
    &protobuf_c_empty_string,
    0,             /* flags */
    0,NULL,NULL    /* reserved1,reserved2, etc */
  },
  {
    "group",
    4,
    PROTOBUF_C_LABEL_NONE,
    PROTOBUF_C_TYPE_STRING,
    0,   /* quantifier_offset */
    offsetof(Auth__Sys, group),
    NULL,
    &protobuf_c_empty_string,
    0,             /* flags */
    0,NULL,NULL    /* reserved1,reserved2, etc */
  },
  {
    "groups",
    5,
    PROTOBUF_C_LABEL_REPEATED,
    PROTOBUF_C_TYPE_STRING,
    offsetof(Auth__Sys, n_groups),
    offsetof(Auth__Sys, groups),
    NULL,
    &protobuf_c_empty_string,
    0,             /* flags */
    0,NULL,NULL    /* reserved1,reserved2, etc */
  },
  {
    "secctx",
    6,
    PROTOBUF_C_LABEL_NONE,
    PROTOBUF_C_TYPE_STRING,
    0,   /* quantifier_offset */
    offsetof(Auth__Sys, secctx),
    NULL,
    &protobuf_c_empty_string,
    0,             /* flags */
    0,NULL,NULL    /* reserved1,reserved2, etc */
  },
};
static const unsigned auth__sys__field_indices_by_name[] = {
  3,   /* field[3] = group */
  4,   /* field[4] = groups */
  1,   /* field[1] = machinename */
  5,   /* field[5] = secctx */
  0,   /* field[0] = stamp */
  2,   /* field[2] = user */
};
static const ProtobufCIntRange auth__sys__number_ranges[1 + 1] =
{
  { 1, 0 },
  { 0, 6 }
};
const ProtobufCMessageDescriptor auth__sys__descriptor =
{
  PROTOBUF_C__MESSAGE_DESCRIPTOR_MAGIC,
  "auth.Sys",
  "Sys",
  "Auth__Sys",
  "auth",
  sizeof(Auth__Sys),
  6,
  auth__sys__field_descriptors,
  auth__sys__field_indices_by_name,
  1,  auth__sys__number_ranges,
  (ProtobufCMessageInit) auth__sys__init,
  NULL,NULL,NULL    /* reserved[123] */
};
static const ProtobufCFieldDescriptor auth__sys_verifier__field_descriptors[1] =
{
  {
    "signature",
    1,
    PROTOBUF_C_LABEL_NONE,
    PROTOBUF_C_TYPE_BYTES,
    0,   /* quantifier_offset */
    offsetof(Auth__SysVerifier, signature),
    NULL,
    NULL,
    0,             /* flags */
    0,NULL,NULL    /* reserved1,reserved2, etc */
  },
};
static const unsigned auth__sys_verifier__field_indices_by_name[] = {
  0,   /* field[0] = signature */
};
static const ProtobufCIntRange auth__sys_verifier__number_ranges[1 + 1] =
{
  { 1, 0 },
  { 0, 1 }
};
const ProtobufCMessageDescriptor auth__sys_verifier__descriptor =
{
  PROTOBUF_C__MESSAGE_DESCRIPTOR_MAGIC,
  "auth.SysVerifier",
  "SysVerifier",
  "Auth__SysVerifier",
  "auth",
  sizeof(Auth__SysVerifier),
  1,
  auth__sys_verifier__field_descriptors,
  auth__sys_verifier__field_indices_by_name,
  1,  auth__sys_verifier__number_ranges,
  (ProtobufCMessageInit) auth__sys_verifier__init,
  NULL,NULL,NULL    /* reserved[123] */
};
static const ProtobufCFieldDescriptor auth__credential__field_descriptors[2] =
{
  {
    "token",
    1,
    PROTOBUF_C_LABEL_NONE,
    PROTOBUF_C_TYPE_MESSAGE,
    0,   /* quantifier_offset */
    offsetof(Auth__Credential, token),
    &auth__token__descriptor,
    NULL,
    0,             /* flags */
    0,NULL,NULL    /* reserved1,reserved2, etc */
  },
  {
    "verifier",
    2,
    PROTOBUF_C_LABEL_NONE,
    PROTOBUF_C_TYPE_MESSAGE,
    0,   /* quantifier_offset */
    offsetof(Auth__Credential, verifier),
    &auth__token__descriptor,
    NULL,
    0,             /* flags */
    0,NULL,NULL    /* reserved1,reserved2, etc */
  },
};
static const unsigned auth__credential__field_indices_by_name[] = {
  0,   /* field[0] = token */
  1,   /* field[1] = verifier */
};
static const ProtobufCIntRange auth__credential__number_ranges[1 + 1] =
{
  { 1, 0 },
  { 0, 2 }
};
const ProtobufCMessageDescriptor auth__credential__descriptor =
{
  PROTOBUF_C__MESSAGE_DESCRIPTOR_MAGIC,
  "auth.Credential",
  "Credential",
  "Auth__Credential",
  "auth",
  sizeof(Auth__Credential),
  2,
  auth__credential__field_descriptors,
  auth__credential__field_indices_by_name,
  1,  auth__credential__number_ranges,
  (ProtobufCMessageInit) auth__credential__init,
  NULL,NULL,NULL    /* reserved[123] */
};
static const ProtobufCEnumValue auth__flavor__enum_values_by_number[2] =
{
  { "AUTH_NONE", "AUTH__FLAVOR__AUTH_NONE", 0 },
  { "AUTH_SYS", "AUTH__FLAVOR__AUTH_SYS", 1 },
};
static const ProtobufCIntRange auth__flavor__value_ranges[] = {
{0, 0},{0, 2}
};
static const ProtobufCEnumValueIndex auth__flavor__enum_values_by_name[2] =
{
  { "AUTH_NONE", 0 },
  { "AUTH_SYS", 1 },
};
const ProtobufCEnumDescriptor auth__flavor__descriptor =
{
  PROTOBUF_C__ENUM_DESCRIPTOR_MAGIC,
  "auth.Flavor",
  "Flavor",
  "Auth__Flavor",
  "auth",
  2,
  auth__flavor__enum_values_by_number,
  2,
  auth__flavor__enum_values_by_name,
  1,
  auth__flavor__value_ranges,
  NULL,NULL,NULL,NULL   /* reserved[1234] */
};
