package io.daos.obj;

import io.daos.BufferAllocator;
import io.daos.Constants;
import io.netty.buffer.ByteBuf;
import org.junit.Assert;
import org.junit.Test;

import java.util.Arrays;

public class IODataDescTest {

  @Test
  public void testKeyLengthZero() throws Exception {
    String dkey;
    int len = Short.MAX_VALUE + 1;
    StringBuilder sb = new StringBuilder();
    for (int i = 0; i < len; i++) {
      sb.append(i);
      if (sb.length() > len) {
        break;
      }
    }
    dkey = sb.toString();
    IllegalArgumentException ee = null;
    try {
      new IODataDesc(dkey, IODataDesc.IodType.ARRAY, 1, true);
    } catch (IllegalArgumentException e) {
      ee = e;
    }
    Assert.assertNotNull(ee);
    Assert.assertTrue(ee.getMessage().contains("should not exceed " + Short.MAX_VALUE));
  }

  @Test
  public void testInconsistentAction() throws Exception {
    IllegalArgumentException ee = null;
    IODataDesc desc = null;
    try {
      desc = new IODataDesc("dkey1", IODataDesc.IodType.ARRAY, 1, true);
      desc.addEntryForFetch("akey1", 0, 10);
    } catch (IllegalArgumentException e) {
      ee = e;
    } finally {
      desc.release();
    }
    Assert.assertNotNull(ee);
    Assert.assertTrue(ee.getMessage().contains("It's desc for update"));
  }

  @Test
  public void testAkeyLengthExceedLimit() throws Exception {
    String akey;
    int len = Short.MAX_VALUE + 1;
    StringBuilder sb = new StringBuilder();
    for (int i = 0; i < len; i++) {
      sb.append(i);
      if (sb.length() > len) {
        break;
      }
    }
    akey = sb.toString();
    IllegalArgumentException ee = null;
    ByteBuf buffer = BufferAllocator.objBufWithNativeOrder(10);
    buffer.writeByte(1);
    IODataDesc desc = null;
    try {
      desc = new IODataDesc("dkey1", IODataDesc.IodType.ARRAY, 1, true);
      desc.addEntryForUpdate(akey, 0, buffer);
    } catch (IllegalArgumentException e) {
      ee = e;
    } finally {
      desc.release();
      buffer.release();
    }
    Assert.assertNotNull(ee);
    Assert.assertTrue(ee.getMessage().contains("should not exceed " + Short.MAX_VALUE));
  }

  @Test
  public void testOffsetNotMultipleOfRecordSize() throws Exception {
    IllegalArgumentException ee = null;
    IODataDesc desc = null;
    try {
      desc = new IODataDesc("dkey1", IODataDesc.IodType.ARRAY, 10, false);
      desc.addEntryForFetch("akey", 9, 10);
    } catch (IllegalArgumentException e) {
      ee = e;
    } finally {
      desc.release();
    }
    Assert.assertNotNull(ee);
    Assert.assertTrue(ee.getMessage().contains("should be a multiple of recordSize"));
  }

  @Test
  public void testNonPositiveDataSize() throws Exception {
    IllegalArgumentException ee = null;
    IODataDesc desc = null;
    try {
      desc = new IODataDesc("dkey1", IODataDesc.IodType.ARRAY, 10, false);
      desc.addEntryForFetch("akey", 10, 0);
    } catch (IllegalArgumentException e) {
      ee = e;
    } finally {
      desc.release();
    }
    Assert.assertNotNull(ee);
    Assert.assertTrue(ee.getMessage().contains("data size should be positive"));
    ee = null;
    try {
      desc = new IODataDesc("dkey1", IODataDesc.IodType.ARRAY, 10, false);
      desc.addEntryForFetch("akey", 10, -1);
    } catch (IllegalArgumentException e) {
      ee = e;
    } finally {
      desc.release();
    }
    Assert.assertNotNull(ee);
    Assert.assertTrue(ee.getMessage().contains("data size should be positive"));
  }

  @Test
  public void testSingleValueNonZeroOffset() throws Exception {
    IllegalArgumentException ee = null;
    IODataDesc desc = null;
    try {
      desc = new IODataDesc("dkey1", IODataDesc.IodType.SINGLE, 10, false);
      desc.addEntryForFetch("akey", 10, 10);
    } catch (IllegalArgumentException e) {
      ee = e;
    } finally {
      desc.release();
    }
    Assert.assertNotNull(ee);
    Assert.assertTrue(ee.getMessage().contains("offset should be zero for"));
  }

  @Test
  public void testSingleValueDataSizeBiggerThanRecSize() throws Exception {
    IllegalArgumentException ee = null;
    IODataDesc desc = null;
    try {
      desc = new IODataDesc("dkey1", IODataDesc.IodType.SINGLE, 10, false);
      desc.addEntryForFetch("akey", 0, 70);
    } catch (IllegalArgumentException e) {
      ee = e;
    } finally {
      desc.release();
    }
    Assert.assertNotNull(ee);
    Assert.assertTrue(ee.getMessage().contains("data size should be no more than record size for"));
  }

  @Test
  public void testInvalidIodType() throws Exception {
    IllegalArgumentException ee = null;
    try {
      new IODataDesc("dkey1", IODataDesc.IodType.NONE, 10, false);
    } catch (IllegalArgumentException e) {
      ee = e;
    }
    Assert.assertNotNull(ee);
    Assert.assertTrue(ee.getMessage().contains("need valid IodType, either"));
  }

  @Test
  public void testCallFetchMethodsWhenUpdate() throws Exception {
    UnsupportedOperationException ee = null;
    ByteBuf buffer = BufferAllocator.objBufWithNativeOrder(10);
    IODataDesc desc = null;
    try {
      buffer.writerIndex(buffer.capacity());
      desc = new IODataDesc("dkey1", IODataDesc.IodType.SINGLE, 10, true);
      IODataDesc.Entry entry = desc.addEntryForUpdate("akey", 0, buffer);
      try {
        entry.getActualSize();
      } catch (UnsupportedOperationException e) {
        ee = e;
      }
      Assert.assertNotNull(ee);
      Assert.assertTrue(ee.getMessage().contains("only support for fetch"));
      ee = null;
      try {
        entry.setActualSize(10);
      } catch (UnsupportedOperationException e) {
        ee = e;
      }
      Assert.assertNotNull(ee);
      Assert.assertTrue(ee.getMessage().contains("only support for fetch"));
      ee = null;
      try {
        entry.getActualRecSize();
      } catch (UnsupportedOperationException e) {
        ee = e;
      }
      Assert.assertNotNull(ee);
      Assert.assertTrue(ee.getMessage().contains("only support for fetch"));
      ee = null;
      try {
        entry.setActualRecSize(10);
      } catch (UnsupportedOperationException e) {
        ee = e;
      }
      Assert.assertNotNull(ee);
      Assert.assertTrue(ee.getMessage().contains("only support for fetch"));
      ee = null;
      try {
        entry.getFetchedData();
      } catch (UnsupportedOperationException e) {
        ee = e;
      }
      Assert.assertNotNull(ee);
      Assert.assertTrue(ee.getMessage().contains("only support for fetch"));
    } finally {
      buffer.release();
      desc.release();
    }
  }

  @Test
  public void testEncodeReusableDesc() throws Exception {
    // array value
//    ByteBuf buffer  = BufferAllocator.objBufWithNativeOrder(30);
//    buffer.writerIndex(30);
    IODataDesc desc = null;
    try {
      desc = new IODataDesc(10, 2, 100, IODataDesc.IodType.ARRAY,
        10, true);
      checkEncoded(desc, IODataDesc.IodType.ARRAY, true);
      // single value
      desc = new IODataDesc(10, 2, 100, IODataDesc.IodType.SINGLE,
        10, false);
      checkEncoded(desc, IODataDesc.IodType.SINGLE, false);
    } finally {
//      buffer.release();
      desc.release();
    }
  }

  private void checkEncoded(IODataDesc desc, IODataDesc.IodType type, boolean update) throws Exception {
    desc.setDkey("dkey");
    for (int i = 0; i < 2; i++) {
      IODataDesc.Entry e = desc.getEntry(i);
      if (update) {
        e.getDataBuffer().writerIndex(30);
        e.setKey("key" + i, 0, e.getDataBuffer());
      } else {
        e.setKey("key" + i, 0, 10);
      }
    }
    desc.encode();
    ByteBuf descBuffer = desc.getDescBuffer();
    if (update) {
      Assert.assertEquals(85, descBuffer.writerIndex());
    } else {
      Assert.assertEquals(69, descBuffer.writerIndex());
      Assert.assertEquals(85, descBuffer.capacity());
    }
    // address
    descBuffer.readerIndex(0);
    Assert.assertEquals(0L, descBuffer.readLong());
    // max key len
    Assert.assertEquals(10, descBuffer.readShort());
    // nbr of akeys with data
    Assert.assertEquals(2, descBuffer.readShort());
    // type
    Assert.assertEquals(type.getValue(), descBuffer.readByte());
    // record size
    Assert.assertEquals(10, descBuffer.readInt());
    // dkey
    Assert.assertEquals(4, descBuffer.readShort());
    byte keyBytes[] = new byte[4];
    descBuffer.readBytes(keyBytes);
    Assert.assertTrue(Arrays.equals("dkey".getBytes(Constants.KEY_CHARSET), keyBytes));
    descBuffer.readerIndex(descBuffer.readerIndex() + 10 - 4);
    // entries
    for (int i = 0; i < 2; i++) {
      Assert.assertEquals(4, descBuffer.readShort());
      descBuffer.readBytes(keyBytes);
      Assert.assertTrue(Arrays.equals(("key" + i).getBytes(Constants.KEY_CHARSET), keyBytes));
      descBuffer.readerIndex(descBuffer.readerIndex() + 10 - 4);
      if (type == IODataDesc.IodType.ARRAY) {
        Assert.assertEquals(0, descBuffer.readInt());
        Assert.assertEquals(3, descBuffer.readInt());
      }
      descBuffer.readerIndex(descBuffer.readerIndex() + 8);
    }
    desc.release();
  }

  @Test
  public void testGetDataWhenFetch() throws Exception {
    // single value
    IODataDesc desc = null;
    try {
      desc = new IODataDesc("dkey", IODataDesc.IodType.SINGLE, 10, false);
      IODataDesc.Entry entry = desc.addEntryForFetch("akey", 0, 10);
      desc.encode();
      Assert.assertEquals(33, desc.getDescBuffer().writerIndex());
      ByteBuf descBuf = desc.getDescBuffer();
      descBuf.writeInt(8);
      descBuf.writeInt(8);
      desc.parseResult();
      // read to byte array
      ByteBuf buf = entry.getFetchedData();
      Assert.assertEquals(0, buf.readerIndex());
      Assert.assertEquals(8, buf.writerIndex());
    } finally {
      if (desc != null) {
        desc.release();
      }
    }
    // array value
    try {
      desc = new IODataDesc("dkey", IODataDesc.IodType.ARRAY, 10, false);
      IODataDesc.Entry entry = desc.addEntryForFetch("akey", 0, 30);
      desc.encode();
      ByteBuf descBuf = desc.getDescBuffer();
      Assert.assertEquals(41, descBuf.writerIndex());
      // not reusable
      descBuf.readerIndex(0);
      Assert.assertEquals(-1L, descBuf.readLong());
      // check iod type and record size
      Assert.assertEquals(IODataDesc.IodType.ARRAY.getValue(), descBuf.readByte());
      Assert.assertEquals(10, descBuf.readInt());
      // check dkey
      Assert.assertEquals(4, descBuf.readShort());
      byte keyBytes[] = new byte[4];
      descBuf.readBytes(keyBytes);
      Assert.assertTrue(Arrays.equals("dkey".getBytes(Constants.KEY_CHARSET), keyBytes));
      Assert.assertEquals(4, descBuf.readShort());
      descBuf.readBytes(keyBytes);
      Assert.assertTrue(Arrays.equals("akey".getBytes(Constants.KEY_CHARSET), keyBytes));
      Assert.assertEquals(0, descBuf.readInt());
      Assert.assertEquals(3, descBuf.readInt());
      // parse
      descBuf.writeInt(30);
      descBuf.writeInt(10);
      desc.parseResult();
      ByteBuf buf = entry.getFetchedData();
      Assert.assertEquals(0, buf.readerIndex());
      Assert.assertEquals(30, buf.writerIndex());
    } finally {
      if (desc != null) {
        desc.release();
      }
    }
  }
}
