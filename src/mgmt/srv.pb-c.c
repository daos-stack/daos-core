/* Generated by the protocol buffer compiler.  DO NOT EDIT! */
/* Generated from: srv.proto */

/* Do not generate deprecated warnings for self */
#ifndef PROTOBUF_C__NO_DEPRECATED
#define PROTOBUF_C__NO_DEPRECATED
#endif

#include "srv.pb-c.h"
void   mgmt__daos_rank__init
                     (Mgmt__DaosRank         *message)
{
  static const Mgmt__DaosRank init_value = MGMT__DAOS_RANK__INIT;
  *message = init_value;
}
size_t mgmt__daos_rank__get_packed_size
                     (const Mgmt__DaosRank *message)
{
  assert(message->base.descriptor == &mgmt__daos_rank__descriptor);
  return protobuf_c_message_get_packed_size ((const ProtobufCMessage*)(message));
}
size_t mgmt__daos_rank__pack
                     (const Mgmt__DaosRank *message,
                      uint8_t       *out)
{
  assert(message->base.descriptor == &mgmt__daos_rank__descriptor);
  return protobuf_c_message_pack ((const ProtobufCMessage*)message, out);
}
size_t mgmt__daos_rank__pack_to_buffer
                     (const Mgmt__DaosRank *message,
                      ProtobufCBuffer *buffer)
{
  assert(message->base.descriptor == &mgmt__daos_rank__descriptor);
  return protobuf_c_message_pack_to_buffer ((const ProtobufCMessage*)message, buffer);
}
Mgmt__DaosRank *
       mgmt__daos_rank__unpack
                     (ProtobufCAllocator  *allocator,
                      size_t               len,
                      const uint8_t       *data)
{
  return (Mgmt__DaosRank *)
     protobuf_c_message_unpack (&mgmt__daos_rank__descriptor,
                                allocator, len, data);
}
void   mgmt__daos_rank__free_unpacked
                     (Mgmt__DaosRank *message,
                      ProtobufCAllocator *allocator)
{
  if(!message)
    return;
  assert(message->base.descriptor == &mgmt__daos_rank__descriptor);
  protobuf_c_message_free_unpacked ((ProtobufCMessage*)message, allocator);
}
void   mgmt__daos_resp__init
                     (Mgmt__DaosResp         *message)
{
  static const Mgmt__DaosResp init_value = MGMT__DAOS_RESP__INIT;
  *message = init_value;
}
size_t mgmt__daos_resp__get_packed_size
                     (const Mgmt__DaosResp *message)
{
  assert(message->base.descriptor == &mgmt__daos_resp__descriptor);
  return protobuf_c_message_get_packed_size ((const ProtobufCMessage*)(message));
}
size_t mgmt__daos_resp__pack
                     (const Mgmt__DaosResp *message,
                      uint8_t       *out)
{
  assert(message->base.descriptor == &mgmt__daos_resp__descriptor);
  return protobuf_c_message_pack ((const ProtobufCMessage*)message, out);
}
size_t mgmt__daos_resp__pack_to_buffer
                     (const Mgmt__DaosResp *message,
                      ProtobufCBuffer *buffer)
{
  assert(message->base.descriptor == &mgmt__daos_resp__descriptor);
  return protobuf_c_message_pack_to_buffer ((const ProtobufCMessage*)message, buffer);
}
Mgmt__DaosResp *
       mgmt__daos_resp__unpack
                     (ProtobufCAllocator  *allocator,
                      size_t               len,
                      const uint8_t       *data)
{
  return (Mgmt__DaosResp *)
     protobuf_c_message_unpack (&mgmt__daos_resp__descriptor,
                                allocator, len, data);
}
void   mgmt__daos_resp__free_unpacked
                     (Mgmt__DaosResp *message,
                      ProtobufCAllocator *allocator)
{
  if(!message)
    return;
  assert(message->base.descriptor == &mgmt__daos_resp__descriptor);
  protobuf_c_message_free_unpacked ((ProtobufCMessage*)message, allocator);
}
void   mgmt__set_rank_req__init
                     (Mgmt__SetRankReq         *message)
{
  static const Mgmt__SetRankReq init_value = MGMT__SET_RANK_REQ__INIT;
  *message = init_value;
}
size_t mgmt__set_rank_req__get_packed_size
                     (const Mgmt__SetRankReq *message)
{
  assert(message->base.descriptor == &mgmt__set_rank_req__descriptor);
  return protobuf_c_message_get_packed_size ((const ProtobufCMessage*)(message));
}
size_t mgmt__set_rank_req__pack
                     (const Mgmt__SetRankReq *message,
                      uint8_t       *out)
{
  assert(message->base.descriptor == &mgmt__set_rank_req__descriptor);
  return protobuf_c_message_pack ((const ProtobufCMessage*)message, out);
}
size_t mgmt__set_rank_req__pack_to_buffer
                     (const Mgmt__SetRankReq *message,
                      ProtobufCBuffer *buffer)
{
  assert(message->base.descriptor == &mgmt__set_rank_req__descriptor);
  return protobuf_c_message_pack_to_buffer ((const ProtobufCMessage*)message, buffer);
}
Mgmt__SetRankReq *
       mgmt__set_rank_req__unpack
                     (ProtobufCAllocator  *allocator,
                      size_t               len,
                      const uint8_t       *data)
{
  return (Mgmt__SetRankReq *)
     protobuf_c_message_unpack (&mgmt__set_rank_req__descriptor,
                                allocator, len, data);
}
void   mgmt__set_rank_req__free_unpacked
                     (Mgmt__SetRankReq *message,
                      ProtobufCAllocator *allocator)
{
  if(!message)
    return;
  assert(message->base.descriptor == &mgmt__set_rank_req__descriptor);
  protobuf_c_message_free_unpacked ((ProtobufCMessage*)message, allocator);
}
void   mgmt__create_ms_req__init
                     (Mgmt__CreateMsReq         *message)
{
  static const Mgmt__CreateMsReq init_value = MGMT__CREATE_MS_REQ__INIT;
  *message = init_value;
}
size_t mgmt__create_ms_req__get_packed_size
                     (const Mgmt__CreateMsReq *message)
{
  assert(message->base.descriptor == &mgmt__create_ms_req__descriptor);
  return protobuf_c_message_get_packed_size ((const ProtobufCMessage*)(message));
}
size_t mgmt__create_ms_req__pack
                     (const Mgmt__CreateMsReq *message,
                      uint8_t       *out)
{
  assert(message->base.descriptor == &mgmt__create_ms_req__descriptor);
  return protobuf_c_message_pack ((const ProtobufCMessage*)message, out);
}
size_t mgmt__create_ms_req__pack_to_buffer
                     (const Mgmt__CreateMsReq *message,
                      ProtobufCBuffer *buffer)
{
  assert(message->base.descriptor == &mgmt__create_ms_req__descriptor);
  return protobuf_c_message_pack_to_buffer ((const ProtobufCMessage*)message, buffer);
}
Mgmt__CreateMsReq *
       mgmt__create_ms_req__unpack
                     (ProtobufCAllocator  *allocator,
                      size_t               len,
                      const uint8_t       *data)
{
  return (Mgmt__CreateMsReq *)
     protobuf_c_message_unpack (&mgmt__create_ms_req__descriptor,
                                allocator, len, data);
}
void   mgmt__create_ms_req__free_unpacked
                     (Mgmt__CreateMsReq *message,
                      ProtobufCAllocator *allocator)
{
  if(!message)
    return;
  assert(message->base.descriptor == &mgmt__create_ms_req__descriptor);
  protobuf_c_message_free_unpacked ((ProtobufCMessage*)message, allocator);
}
static const ProtobufCFieldDescriptor mgmt__daos_rank__field_descriptors[2] =
{
  {
    "pool_uuid",
    1,
    PROTOBUF_C_LABEL_NONE,
    PROTOBUF_C_TYPE_STRING,
    0,   /* quantifier_offset */
    offsetof(Mgmt__DaosRank, pool_uuid),
    NULL,
    &protobuf_c_empty_string,
    0,             /* flags */
    0,NULL,NULL    /* reserved1,reserved2, etc */
  },
  {
    "rank",
    2,
    PROTOBUF_C_LABEL_NONE,
    PROTOBUF_C_TYPE_UINT32,
    0,   /* quantifier_offset */
    offsetof(Mgmt__DaosRank, rank),
    NULL,
    NULL,
    0,             /* flags */
    0,NULL,NULL    /* reserved1,reserved2, etc */
  },
};
static const unsigned mgmt__daos_rank__field_indices_by_name[] = {
  0,   /* field[0] = pool_uuid */
  1,   /* field[1] = rank */
};
static const ProtobufCIntRange mgmt__daos_rank__number_ranges[1 + 1] =
{
  { 1, 0 },
  { 0, 2 }
};
const ProtobufCMessageDescriptor mgmt__daos_rank__descriptor =
{
  PROTOBUF_C__MESSAGE_DESCRIPTOR_MAGIC,
  "mgmt.DaosRank",
  "DaosRank",
  "Mgmt__DaosRank",
  "mgmt",
  sizeof(Mgmt__DaosRank),
  2,
  mgmt__daos_rank__field_descriptors,
  mgmt__daos_rank__field_indices_by_name,
  1,  mgmt__daos_rank__number_ranges,
  (ProtobufCMessageInit) mgmt__daos_rank__init,
  NULL,NULL,NULL    /* reserved[123] */
};
static const ProtobufCFieldDescriptor mgmt__daos_resp__field_descriptors[1] =
{
  {
    "status",
    1,
    PROTOBUF_C_LABEL_NONE,
    PROTOBUF_C_TYPE_INT32,
    0,   /* quantifier_offset */
    offsetof(Mgmt__DaosResp, status),
    NULL,
    NULL,
    0,             /* flags */
    0,NULL,NULL    /* reserved1,reserved2, etc */
  },
};
static const unsigned mgmt__daos_resp__field_indices_by_name[] = {
  0,   /* field[0] = status */
};
static const ProtobufCIntRange mgmt__daos_resp__number_ranges[1 + 1] =
{
  { 1, 0 },
  { 0, 1 }
};
const ProtobufCMessageDescriptor mgmt__daos_resp__descriptor =
{
  PROTOBUF_C__MESSAGE_DESCRIPTOR_MAGIC,
  "mgmt.DaosResp",
  "DaosResp",
  "Mgmt__DaosResp",
  "mgmt",
  sizeof(Mgmt__DaosResp),
  1,
  mgmt__daos_resp__field_descriptors,
  mgmt__daos_resp__field_indices_by_name,
  1,  mgmt__daos_resp__number_ranges,
  (ProtobufCMessageInit) mgmt__daos_resp__init,
  NULL,NULL,NULL    /* reserved[123] */
};
static const ProtobufCFieldDescriptor mgmt__set_rank_req__field_descriptors[1] =
{
  {
    "rank",
    1,
    PROTOBUF_C_LABEL_NONE,
    PROTOBUF_C_TYPE_UINT32,
    0,   /* quantifier_offset */
    offsetof(Mgmt__SetRankReq, rank),
    NULL,
    NULL,
    0,             /* flags */
    0,NULL,NULL    /* reserved1,reserved2, etc */
  },
};
static const unsigned mgmt__set_rank_req__field_indices_by_name[] = {
  0,   /* field[0] = rank */
};
static const ProtobufCIntRange mgmt__set_rank_req__number_ranges[1 + 1] =
{
  { 1, 0 },
  { 0, 1 }
};
const ProtobufCMessageDescriptor mgmt__set_rank_req__descriptor =
{
  PROTOBUF_C__MESSAGE_DESCRIPTOR_MAGIC,
  "mgmt.SetRankReq",
  "SetRankReq",
  "Mgmt__SetRankReq",
  "mgmt",
  sizeof(Mgmt__SetRankReq),
  1,
  mgmt__set_rank_req__field_descriptors,
  mgmt__set_rank_req__field_indices_by_name,
  1,  mgmt__set_rank_req__number_ranges,
  (ProtobufCMessageInit) mgmt__set_rank_req__init,
  NULL,NULL,NULL    /* reserved[123] */
};
static const ProtobufCFieldDescriptor mgmt__create_ms_req__field_descriptors[3] =
{
  {
    "bootstrap",
    1,
    PROTOBUF_C_LABEL_NONE,
    PROTOBUF_C_TYPE_BOOL,
    0,   /* quantifier_offset */
    offsetof(Mgmt__CreateMsReq, bootstrap),
    NULL,
    NULL,
    0,             /* flags */
    0,NULL,NULL    /* reserved1,reserved2, etc */
  },
  {
    "uuid",
    2,
    PROTOBUF_C_LABEL_NONE,
    PROTOBUF_C_TYPE_STRING,
    0,   /* quantifier_offset */
    offsetof(Mgmt__CreateMsReq, uuid),
    NULL,
    &protobuf_c_empty_string,
    0,             /* flags */
    0,NULL,NULL    /* reserved1,reserved2, etc */
  },
  {
    "addr",
    3,
    PROTOBUF_C_LABEL_NONE,
    PROTOBUF_C_TYPE_STRING,
    0,   /* quantifier_offset */
    offsetof(Mgmt__CreateMsReq, addr),
    NULL,
    &protobuf_c_empty_string,
    0,             /* flags */
    0,NULL,NULL    /* reserved1,reserved2, etc */
  },
};
static const unsigned mgmt__create_ms_req__field_indices_by_name[] = {
  2,   /* field[2] = addr */
  0,   /* field[0] = bootstrap */
  1,   /* field[1] = uuid */
};
static const ProtobufCIntRange mgmt__create_ms_req__number_ranges[1 + 1] =
{
  { 1, 0 },
  { 0, 3 }
};
const ProtobufCMessageDescriptor mgmt__create_ms_req__descriptor =
{
  PROTOBUF_C__MESSAGE_DESCRIPTOR_MAGIC,
  "mgmt.CreateMsReq",
  "CreateMsReq",
  "Mgmt__CreateMsReq",
  "mgmt",
  sizeof(Mgmt__CreateMsReq),
  3,
  mgmt__create_ms_req__field_descriptors,
  mgmt__create_ms_req__field_indices_by_name,
  1,  mgmt__create_ms_req__number_ranges,
  (ProtobufCMessageInit) mgmt__create_ms_req__init,
  NULL,NULL,NULL    /* reserved[123] */
};
