/* Generated by the protocol buffer compiler.  DO NOT EDIT! */
/* Generated from: storage_query.proto */

/* Do not generate deprecated warnings for self */
#ifndef PROTOBUF_C__NO_DEPRECATED
#define PROTOBUF_C__NO_DEPRECATED
#endif

#include "storage_query.pb-c.h"
void   mgmt__bio_health_req__init
                     (Mgmt__BioHealthReq         *message)
{
  static const Mgmt__BioHealthReq init_value = MGMT__BIO_HEALTH_REQ__INIT;
  *message = init_value;
}
size_t mgmt__bio_health_req__get_packed_size
                     (const Mgmt__BioHealthReq *message)
{
  assert(message->base.descriptor == &mgmt__bio_health_req__descriptor);
  return protobuf_c_message_get_packed_size ((const ProtobufCMessage*)(message));
}
size_t mgmt__bio_health_req__pack
                     (const Mgmt__BioHealthReq *message,
                      uint8_t       *out)
{
  assert(message->base.descriptor == &mgmt__bio_health_req__descriptor);
  return protobuf_c_message_pack ((const ProtobufCMessage*)message, out);
}
size_t mgmt__bio_health_req__pack_to_buffer
                     (const Mgmt__BioHealthReq *message,
                      ProtobufCBuffer *buffer)
{
  assert(message->base.descriptor == &mgmt__bio_health_req__descriptor);
  return protobuf_c_message_pack_to_buffer ((const ProtobufCMessage*)message, buffer);
}
Mgmt__BioHealthReq *
       mgmt__bio_health_req__unpack
                     (ProtobufCAllocator  *allocator,
                      size_t               len,
                      const uint8_t       *data)
{
  return (Mgmt__BioHealthReq *)
     protobuf_c_message_unpack (&mgmt__bio_health_req__descriptor,
                                allocator, len, data);
}
void   mgmt__bio_health_req__free_unpacked
                     (Mgmt__BioHealthReq *message,
                      ProtobufCAllocator *allocator)
{
  if(!message)
    return;
  assert(message->base.descriptor == &mgmt__bio_health_req__descriptor);
  protobuf_c_message_free_unpacked ((ProtobufCMessage*)message, allocator);
}
void   mgmt__bio_health_resp__init
                     (Mgmt__BioHealthResp         *message)
{
  static const Mgmt__BioHealthResp init_value = MGMT__BIO_HEALTH_RESP__INIT;
  *message = init_value;
}
size_t mgmt__bio_health_resp__get_packed_size
                     (const Mgmt__BioHealthResp *message)
{
  assert(message->base.descriptor == &mgmt__bio_health_resp__descriptor);
  return protobuf_c_message_get_packed_size ((const ProtobufCMessage*)(message));
}
size_t mgmt__bio_health_resp__pack
                     (const Mgmt__BioHealthResp *message,
                      uint8_t       *out)
{
  assert(message->base.descriptor == &mgmt__bio_health_resp__descriptor);
  return protobuf_c_message_pack ((const ProtobufCMessage*)message, out);
}
size_t mgmt__bio_health_resp__pack_to_buffer
                     (const Mgmt__BioHealthResp *message,
                      ProtobufCBuffer *buffer)
{
  assert(message->base.descriptor == &mgmt__bio_health_resp__descriptor);
  return protobuf_c_message_pack_to_buffer ((const ProtobufCMessage*)message, buffer);
}
Mgmt__BioHealthResp *
       mgmt__bio_health_resp__unpack
                     (ProtobufCAllocator  *allocator,
                      size_t               len,
                      const uint8_t       *data)
{
  return (Mgmt__BioHealthResp *)
     protobuf_c_message_unpack (&mgmt__bio_health_resp__descriptor,
                                allocator, len, data);
}
void   mgmt__bio_health_resp__free_unpacked
                     (Mgmt__BioHealthResp *message,
                      ProtobufCAllocator *allocator)
{
  if(!message)
    return;
  assert(message->base.descriptor == &mgmt__bio_health_resp__descriptor);
  protobuf_c_message_free_unpacked ((ProtobufCMessage*)message, allocator);
}
void   mgmt__smd_dev_req__init
                     (Mgmt__SmdDevReq         *message)
{
  static const Mgmt__SmdDevReq init_value = MGMT__SMD_DEV_REQ__INIT;
  *message = init_value;
}
size_t mgmt__smd_dev_req__get_packed_size
                     (const Mgmt__SmdDevReq *message)
{
  assert(message->base.descriptor == &mgmt__smd_dev_req__descriptor);
  return protobuf_c_message_get_packed_size ((const ProtobufCMessage*)(message));
}
size_t mgmt__smd_dev_req__pack
                     (const Mgmt__SmdDevReq *message,
                      uint8_t       *out)
{
  assert(message->base.descriptor == &mgmt__smd_dev_req__descriptor);
  return protobuf_c_message_pack ((const ProtobufCMessage*)message, out);
}
size_t mgmt__smd_dev_req__pack_to_buffer
                     (const Mgmt__SmdDevReq *message,
                      ProtobufCBuffer *buffer)
{
  assert(message->base.descriptor == &mgmt__smd_dev_req__descriptor);
  return protobuf_c_message_pack_to_buffer ((const ProtobufCMessage*)message, buffer);
}
Mgmt__SmdDevReq *
       mgmt__smd_dev_req__unpack
                     (ProtobufCAllocator  *allocator,
                      size_t               len,
                      const uint8_t       *data)
{
  return (Mgmt__SmdDevReq *)
     protobuf_c_message_unpack (&mgmt__smd_dev_req__descriptor,
                                allocator, len, data);
}
void   mgmt__smd_dev_req__free_unpacked
                     (Mgmt__SmdDevReq *message,
                      ProtobufCAllocator *allocator)
{
  if(!message)
    return;
  assert(message->base.descriptor == &mgmt__smd_dev_req__descriptor);
  protobuf_c_message_free_unpacked ((ProtobufCMessage*)message, allocator);
}
void   mgmt__smd_dev_resp__device__init
                     (Mgmt__SmdDevResp__Device         *message)
{
  static const Mgmt__SmdDevResp__Device init_value = MGMT__SMD_DEV_RESP__DEVICE__INIT;
  *message = init_value;
}
void   mgmt__smd_dev_resp__init
                     (Mgmt__SmdDevResp         *message)
{
  static const Mgmt__SmdDevResp init_value = MGMT__SMD_DEV_RESP__INIT;
  *message = init_value;
}
size_t mgmt__smd_dev_resp__get_packed_size
                     (const Mgmt__SmdDevResp *message)
{
  assert(message->base.descriptor == &mgmt__smd_dev_resp__descriptor);
  return protobuf_c_message_get_packed_size ((const ProtobufCMessage*)(message));
}
size_t mgmt__smd_dev_resp__pack
                     (const Mgmt__SmdDevResp *message,
                      uint8_t       *out)
{
  assert(message->base.descriptor == &mgmt__smd_dev_resp__descriptor);
  return protobuf_c_message_pack ((const ProtobufCMessage*)message, out);
}
size_t mgmt__smd_dev_resp__pack_to_buffer
                     (const Mgmt__SmdDevResp *message,
                      ProtobufCBuffer *buffer)
{
  assert(message->base.descriptor == &mgmt__smd_dev_resp__descriptor);
  return protobuf_c_message_pack_to_buffer ((const ProtobufCMessage*)message, buffer);
}
Mgmt__SmdDevResp *
       mgmt__smd_dev_resp__unpack
                     (ProtobufCAllocator  *allocator,
                      size_t               len,
                      const uint8_t       *data)
{
  return (Mgmt__SmdDevResp *)
     protobuf_c_message_unpack (&mgmt__smd_dev_resp__descriptor,
                                allocator, len, data);
}
void   mgmt__smd_dev_resp__free_unpacked
                     (Mgmt__SmdDevResp *message,
                      ProtobufCAllocator *allocator)
{
  if(!message)
    return;
  assert(message->base.descriptor == &mgmt__smd_dev_resp__descriptor);
  protobuf_c_message_free_unpacked ((ProtobufCMessage*)message, allocator);
}
static const ProtobufCFieldDescriptor mgmt__bio_health_req__field_descriptors[2] =
{
  {
    "dev_uuid",
    1,
    PROTOBUF_C_LABEL_NONE,
    PROTOBUF_C_TYPE_STRING,
    0,   /* quantifier_offset */
    offsetof(Mgmt__BioHealthReq, dev_uuid),
    NULL,
    &protobuf_c_empty_string,
    0,             /* flags */
    0,NULL,NULL    /* reserved1,reserved2, etc */
  },
  {
    "tgt_id",
    2,
    PROTOBUF_C_LABEL_NONE,
    PROTOBUF_C_TYPE_STRING,
    0,   /* quantifier_offset */
    offsetof(Mgmt__BioHealthReq, tgt_id),
    NULL,
    &protobuf_c_empty_string,
    0,             /* flags */
    0,NULL,NULL    /* reserved1,reserved2, etc */
  },
};
static const unsigned mgmt__bio_health_req__field_indices_by_name[] = {
  0,   /* field[0] = dev_uuid */
  1,   /* field[1] = tgt_id */
};
static const ProtobufCIntRange mgmt__bio_health_req__number_ranges[1 + 1] =
{
  { 1, 0 },
  { 0, 2 }
};
const ProtobufCMessageDescriptor mgmt__bio_health_req__descriptor =
{
  PROTOBUF_C__MESSAGE_DESCRIPTOR_MAGIC,
  "mgmt.BioHealthReq",
  "BioHealthReq",
  "Mgmt__BioHealthReq",
  "mgmt",
  sizeof(Mgmt__BioHealthReq),
  2,
  mgmt__bio_health_req__field_descriptors,
  mgmt__bio_health_req__field_indices_by_name,
  1,  mgmt__bio_health_req__number_ranges,
  (ProtobufCMessageInit) mgmt__bio_health_req__init,
  NULL,NULL,NULL    /* reserved[123] */
};
static const ProtobufCFieldDescriptor mgmt__bio_health_resp__field_descriptors[13] =
{
  {
    "status",
    1,
    PROTOBUF_C_LABEL_NONE,
    PROTOBUF_C_TYPE_INT32,
    0,   /* quantifier_offset */
    offsetof(Mgmt__BioHealthResp, status),
    NULL,
    NULL,
    0,             /* flags */
    0,NULL,NULL    /* reserved1,reserved2, etc */
  },
  {
    "dev_uuid",
    2,
    PROTOBUF_C_LABEL_NONE,
    PROTOBUF_C_TYPE_STRING,
    0,   /* quantifier_offset */
    offsetof(Mgmt__BioHealthResp, dev_uuid),
    NULL,
    &protobuf_c_empty_string,
    0,             /* flags */
    0,NULL,NULL    /* reserved1,reserved2, etc */
  },
  {
    "error_count",
    3,
    PROTOBUF_C_LABEL_NONE,
    PROTOBUF_C_TYPE_UINT64,
    0,   /* quantifier_offset */
    offsetof(Mgmt__BioHealthResp, error_count),
    NULL,
    NULL,
    0,             /* flags */
    0,NULL,NULL    /* reserved1,reserved2, etc */
  },
  {
    "temperature",
    4,
    PROTOBUF_C_LABEL_NONE,
    PROTOBUF_C_TYPE_UINT32,
    0,   /* quantifier_offset */
    offsetof(Mgmt__BioHealthResp, temperature),
    NULL,
    NULL,
    0,             /* flags */
    0,NULL,NULL    /* reserved1,reserved2, etc */
  },
  {
    "media_errors",
    5,
    PROTOBUF_C_LABEL_NONE,
    PROTOBUF_C_TYPE_UINT64,
    0,   /* quantifier_offset */
    offsetof(Mgmt__BioHealthResp, media_errors),
    NULL,
    NULL,
    0,             /* flags */
    0,NULL,NULL    /* reserved1,reserved2, etc */
  },
  {
    "read_errs",
    6,
    PROTOBUF_C_LABEL_NONE,
    PROTOBUF_C_TYPE_UINT32,
    0,   /* quantifier_offset */
    offsetof(Mgmt__BioHealthResp, read_errs),
    NULL,
    NULL,
    0,             /* flags */
    0,NULL,NULL    /* reserved1,reserved2, etc */
  },
  {
    "write_errs",
    7,
    PROTOBUF_C_LABEL_NONE,
    PROTOBUF_C_TYPE_UINT32,
    0,   /* quantifier_offset */
    offsetof(Mgmt__BioHealthResp, write_errs),
    NULL,
    NULL,
    0,             /* flags */
    0,NULL,NULL    /* reserved1,reserved2, etc */
  },
  {
    "unmap_errs",
    8,
    PROTOBUF_C_LABEL_NONE,
    PROTOBUF_C_TYPE_UINT32,
    0,   /* quantifier_offset */
    offsetof(Mgmt__BioHealthResp, unmap_errs),
    NULL,
    NULL,
    0,             /* flags */
    0,NULL,NULL    /* reserved1,reserved2, etc */
  },
  {
    "temp",
    9,
    PROTOBUF_C_LABEL_NONE,
    PROTOBUF_C_TYPE_BOOL,
    0,   /* quantifier_offset */
    offsetof(Mgmt__BioHealthResp, temp),
    NULL,
    NULL,
    0,             /* flags */
    0,NULL,NULL    /* reserved1,reserved2, etc */
  },
  {
    "spare",
    10,
    PROTOBUF_C_LABEL_NONE,
    PROTOBUF_C_TYPE_BOOL,
    0,   /* quantifier_offset */
    offsetof(Mgmt__BioHealthResp, spare),
    NULL,
    NULL,
    0,             /* flags */
    0,NULL,NULL    /* reserved1,reserved2, etc */
  },
  {
    "readonly",
    11,
    PROTOBUF_C_LABEL_NONE,
    PROTOBUF_C_TYPE_BOOL,
    0,   /* quantifier_offset */
    offsetof(Mgmt__BioHealthResp, readonly),
    NULL,
    NULL,
    0,             /* flags */
    0,NULL,NULL    /* reserved1,reserved2, etc */
  },
  {
    "device_reliability",
    12,
    PROTOBUF_C_LABEL_NONE,
    PROTOBUF_C_TYPE_BOOL,
    0,   /* quantifier_offset */
    offsetof(Mgmt__BioHealthResp, device_reliability),
    NULL,
    NULL,
    0,             /* flags */
    0,NULL,NULL    /* reserved1,reserved2, etc */
  },
  {
    "volatile_memory",
    13,
    PROTOBUF_C_LABEL_NONE,
    PROTOBUF_C_TYPE_BOOL,
    0,   /* quantifier_offset */
    offsetof(Mgmt__BioHealthResp, volatile_memory),
    NULL,
    NULL,
    0,             /* flags */
    0,NULL,NULL    /* reserved1,reserved2, etc */
  },
};
static const unsigned mgmt__bio_health_resp__field_indices_by_name[] = {
  1,   /* field[1] = dev_uuid */
  11,   /* field[11] = device_reliability */
  2,   /* field[2] = error_count */
  4,   /* field[4] = media_errors */
  5,   /* field[5] = read_errs */
  10,   /* field[10] = readonly */
  9,   /* field[9] = spare */
  0,   /* field[0] = status */
  8,   /* field[8] = temp */
  3,   /* field[3] = temperature */
  7,   /* field[7] = unmap_errs */
  12,   /* field[12] = volatile_memory */
  6,   /* field[6] = write_errs */
};
static const ProtobufCIntRange mgmt__bio_health_resp__number_ranges[1 + 1] =
{
  { 1, 0 },
  { 0, 13 }
};
const ProtobufCMessageDescriptor mgmt__bio_health_resp__descriptor =
{
  PROTOBUF_C__MESSAGE_DESCRIPTOR_MAGIC,
  "mgmt.BioHealthResp",
  "BioHealthResp",
  "Mgmt__BioHealthResp",
  "mgmt",
  sizeof(Mgmt__BioHealthResp),
  13,
  mgmt__bio_health_resp__field_descriptors,
  mgmt__bio_health_resp__field_indices_by_name,
  1,  mgmt__bio_health_resp__number_ranges,
  (ProtobufCMessageInit) mgmt__bio_health_resp__init,
  NULL,NULL,NULL    /* reserved[123] */
};
#define mgmt__smd_dev_req__field_descriptors NULL
#define mgmt__smd_dev_req__field_indices_by_name NULL
#define mgmt__smd_dev_req__number_ranges NULL
const ProtobufCMessageDescriptor mgmt__smd_dev_req__descriptor =
{
  PROTOBUF_C__MESSAGE_DESCRIPTOR_MAGIC,
  "mgmt.SmdDevReq",
  "SmdDevReq",
  "Mgmt__SmdDevReq",
  "mgmt",
  sizeof(Mgmt__SmdDevReq),
  0,
  mgmt__smd_dev_req__field_descriptors,
  mgmt__smd_dev_req__field_indices_by_name,
  0,  mgmt__smd_dev_req__number_ranges,
  (ProtobufCMessageInit) mgmt__smd_dev_req__init,
  NULL,NULL,NULL    /* reserved[123] */
};
static const ProtobufCFieldDescriptor mgmt__smd_dev_resp__device__field_descriptors[2] =
{
  {
    "uuid",
    1,
    PROTOBUF_C_LABEL_NONE,
    PROTOBUF_C_TYPE_STRING,
    0,   /* quantifier_offset */
    offsetof(Mgmt__SmdDevResp__Device, uuid),
    NULL,
    &protobuf_c_empty_string,
    0,             /* flags */
    0,NULL,NULL    /* reserved1,reserved2, etc */
  },
  {
    "tgt_ids",
    2,
    PROTOBUF_C_LABEL_REPEATED,
    PROTOBUF_C_TYPE_INT32,
    offsetof(Mgmt__SmdDevResp__Device, n_tgt_ids),
    offsetof(Mgmt__SmdDevResp__Device, tgt_ids),
    NULL,
    NULL,
    0 | PROTOBUF_C_FIELD_FLAG_PACKED,             /* flags */
    0,NULL,NULL    /* reserved1,reserved2, etc */
  },
};
static const unsigned mgmt__smd_dev_resp__device__field_indices_by_name[] = {
  1,   /* field[1] = tgt_ids */
  0,   /* field[0] = uuid */
};
static const ProtobufCIntRange mgmt__smd_dev_resp__device__number_ranges[1 + 1] =
{
  { 1, 0 },
  { 0, 2 }
};
const ProtobufCMessageDescriptor mgmt__smd_dev_resp__device__descriptor =
{
  PROTOBUF_C__MESSAGE_DESCRIPTOR_MAGIC,
  "mgmt.SmdDevResp.Device",
  "Device",
  "Mgmt__SmdDevResp__Device",
  "mgmt",
  sizeof(Mgmt__SmdDevResp__Device),
  2,
  mgmt__smd_dev_resp__device__field_descriptors,
  mgmt__smd_dev_resp__device__field_indices_by_name,
  1,  mgmt__smd_dev_resp__device__number_ranges,
  (ProtobufCMessageInit) mgmt__smd_dev_resp__device__init,
  NULL,NULL,NULL    /* reserved[123] */
};
static const ProtobufCFieldDescriptor mgmt__smd_dev_resp__field_descriptors[2] =
{
  {
    "status",
    1,
    PROTOBUF_C_LABEL_NONE,
    PROTOBUF_C_TYPE_INT32,
    0,   /* quantifier_offset */
    offsetof(Mgmt__SmdDevResp, status),
    NULL,
    NULL,
    0,             /* flags */
    0,NULL,NULL    /* reserved1,reserved2, etc */
  },
  {
    "devices",
    2,
    PROTOBUF_C_LABEL_REPEATED,
    PROTOBUF_C_TYPE_MESSAGE,
    offsetof(Mgmt__SmdDevResp, n_devices),
    offsetof(Mgmt__SmdDevResp, devices),
    &mgmt__smd_dev_resp__device__descriptor,
    NULL,
    0,             /* flags */
    0,NULL,NULL    /* reserved1,reserved2, etc */
  },
};
static const unsigned mgmt__smd_dev_resp__field_indices_by_name[] = {
  1,   /* field[1] = devices */
  0,   /* field[0] = status */
};
static const ProtobufCIntRange mgmt__smd_dev_resp__number_ranges[1 + 1] =
{
  { 1, 0 },
  { 0, 2 }
};
const ProtobufCMessageDescriptor mgmt__smd_dev_resp__descriptor =
{
  PROTOBUF_C__MESSAGE_DESCRIPTOR_MAGIC,
  "mgmt.SmdDevResp",
  "SmdDevResp",
  "Mgmt__SmdDevResp",
  "mgmt",
  sizeof(Mgmt__SmdDevResp),
  2,
  mgmt__smd_dev_resp__field_descriptors,
  mgmt__smd_dev_resp__field_indices_by_name,
  1,  mgmt__smd_dev_resp__number_ranges,
  (ProtobufCMessageInit) mgmt__smd_dev_resp__init,
  NULL,NULL,NULL    /* reserved[123] */
};
