/*
 * (C) Copyright 2018-2021 Intel Corporation.
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

package io.daos.obj;

import java.io.IOException;

/**
 * Exception with DAOS object info.
 */
public class DaosObjectException extends IOException {

  private DaosObjectId oid;

  public DaosObjectException(DaosObjectId oid, String msg, Throwable cause) {
    super(msg, cause);
    this.oid = oid;
  }

  public DaosObjectException(DaosObjectId oid, String msg) {
    super(msg);
    this.oid = oid;
  }

  /**
   * get message.
   * @return
   */
  @Override
  public String getMessage() {
    return toString();
  }

  /**
   * get localized message.
   * @return
   */
  @Override
  public String getLocalizedMessage() {
    return toString();
  }

  @Override
  public String toString() {
    return super.getMessage() + " " + oid;
  }
}
