package com.intel.daos.client;

import org.junit.Assert;
import org.junit.BeforeClass;
import org.junit.Test;

import java.io.IOException;

public class DaosFsClientIT {

  private static String poolId;
  private static String contId;

  @BeforeClass
  public static void setup()throws Exception{
    poolId = System.getProperty("pool_id", DaosFsClientTestBase.DEFAULT_POOL_ID);
    contId = System.getProperty("cont_id", DaosFsClientTestBase.DEFAULT_CONT_ID);
  }

  @Test
  public void testCreateFsClientFromSpecifiedContainer() throws Exception{
    DaosFsClient.DaosFsClientBuilder builder = new DaosFsClient.DaosFsClientBuilder();
    builder.poolId(poolId).containerId(contId);
    DaosFsClient client = null;
    try{
      client = builder.build();
      Assert.assertTrue(client != null);
    }finally {
      if(client != null){
        client.disconnect();
      }
    }
  }

  @Test
  public void testCreateFsClientFromRootContainer() throws Exception{
    DaosFsClient.DaosFsClientBuilder builder = new DaosFsClient.DaosFsClientBuilder();
    builder.poolId(poolId);
    DaosFsClient client = null;
    try{
      client = builder.build();
      Assert.assertTrue(client != null);
    }finally {
      if(client != null){
        client.disconnect();
      }
    }
  }

  @Test
  public void testCreateNewPoolWithoutScmSize()throws Exception{
    DaosFsClient.DaosFsClientBuilder builder = new DaosFsClient.DaosFsClientBuilder();
    Exception ee = null;
    DaosFsClient client = null;
    try {
      client = builder.build();
    }catch (Exception e){
      ee = e;
    }finally {
      if(client != null){
        client.disconnect();
      }
    }
    Assert.assertTrue(ee instanceof DaosIOException);
    DaosIOException de = (DaosIOException)ee;
    Assert.assertEquals(Constants.CUSTOM_ERR_NO_POOL_SIZE.getCode(), de.getErrorCode());
  }

  @Test
  public void testCreateNewPool()throws Exception {
    DaosFsClient.DaosFsClientBuilder builder = new DaosFsClient.DaosFsClientBuilder();
    builder.poolScmSize(1*1024*1024*1024);
    DaosFsClient client = null;
    try{
      client = builder.build();
      Assert.assertTrue(client != null);
    }finally {
      if(client != null){
        client.disconnect();
        DaosFsClient.destroyPool(Constants.POOL_DEFAULT_SERVER_GROUP, client.getPoolId(), true);
      }
    }
  }

  @Test
  public void testFsClientCachePerPoolAndContainer()throws Exception{
    DaosFsClient.DaosFsClientBuilder builder = new DaosFsClient.DaosFsClientBuilder();
    builder.poolId(poolId).containerId(contId);
    DaosFsClient client = null;
    DaosFsClient client2[] = new DaosFsClient[1];
    try {
      client = builder.build();
      Thread thread = new Thread(){
        @Override
        public void run(){
          try {
            DaosFsClient.DaosFsClientBuilder builder = new DaosFsClient.DaosFsClientBuilder();
            builder.poolId(poolId).containerId(contId);
            client2[0] = builder.build();
          }catch (IOException e){
            e.printStackTrace();
          }
        }
      };
      thread.start();
      thread.join();
      Assert.assertEquals(client, client2[0]);
    } finally {
      client.disconnect();
    }
  }

  @Test
  public void testDeleteSuccessful()throws Exception {
    DaosFsClient.DaosFsClientBuilder builder = new DaosFsClient.DaosFsClientBuilder();
    builder.poolId(poolId).containerId(contId);
    DaosFsClient client = null;
    try{
      client = builder.build();
      Assert.assertTrue(client != null);
      client.delete("/ddddddd/zyx", true);
    }finally {
      if(client != null){
        client.disconnect();
      }
    }
  }

  @Test
  public void testMultipleMkdirs()throws Exception {
    DaosFsClient.DaosFsClientBuilder builder = new DaosFsClient.DaosFsClientBuilder();
    builder.poolId(poolId).containerId(contId);
    DaosFsClient client = null;
    try{
      client = builder.build();
      Assert.assertTrue(client != null);
      client.mkdir("/mkdirs/1", true);
      client.mkdir("/mkdirs/1", true);
    }finally {
      if(client != null){
        client.disconnect();
      }
    }
  }

  @Test(expected = IOException.class)
  public void testMultipleMkdir()throws Exception {
    DaosFsClient.DaosFsClientBuilder builder = new DaosFsClient.DaosFsClientBuilder();
    builder.poolId(poolId).containerId(contId);
    DaosFsClient client = null;
    try{
      client = builder.build();
      Assert.assertTrue(client != null);
      client.mkdir("/mkdir/1", false);
      client.mkdir("/mkdir/1", false);
    }finally {
      if(client != null){
        client.disconnect();
      }
    }
  }

  @Test
  public void testFsClientReferenceOne() throws Exception {
    DaosFsClient.DaosFsClientBuilder builder = new DaosFsClient.DaosFsClientBuilder();
    builder.poolId(poolId).containerId(contId);
    DaosFsClient client = null;
    try{
      client = builder.build();
      Assert.assertEquals(1, client.getRefCnt());
    }finally {
      if(client != null){
        client.disconnect();
        Assert.assertEquals(0, client.getRefCnt());
      }
    }
    Exception ee = null;
    try {
      client.incrementRef();
    } catch (Exception e) {
      ee = e;
    }
    Assert.assertTrue(ee instanceof IllegalStateException);
  }

  @Test
  public void testFsClientReferenceMore() throws Exception {
    DaosFsClient.DaosFsClientBuilder builder = new DaosFsClient.DaosFsClientBuilder();
    builder.poolId(poolId).containerId(contId);
    DaosFsClient client = null;
    int cnt = 0;
    try{
      client = builder.build();
      cnt = client.getRefCnt();
      builder.build();
      Assert.assertEquals(cnt + 1, client.getRefCnt());
      client.disconnect();
      client.incrementRef();
      Assert.assertEquals(cnt + 1, client.getRefCnt());
      client.decrementRef();
    } catch (Exception e){
      e.printStackTrace();
    }finally {
      if(client != null){
        client.disconnect();
        Assert.assertEquals(cnt - 1, client.getRefCnt());
      }
    }
  }
}
