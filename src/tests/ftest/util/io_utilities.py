#!/usr/bin/python
'''
  (C) Copyright 2018-2019 Intel Corporation.

  Licensed under the Apache License, Version 2.0 (the "License");
  you may not use this file except in compliance with the License.
  You may obtain a copy of the License at

     http://www.apache.org/licenses/LICENSE-2.0

  Unless required by applicable law or agreed to in writing, software
  distributed under the License is distributed on an "AS IS" BASIS,
  WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
  See the License for the specific language governing permissions and
  limitations under the License.

  GOVERNMENT LICENSE RIGHTS-OPEN SOURCE SOFTWARE
  The Government's rights to use, modify, reproduce, release, perform, display,
  or disclose this software are subject to the terms of the Apache License as
  provided in Contract No. B609815.
  Any reproduction of computer software, computer software documentation, or
  portions thereof marked with this legend must also reproduce the markings.
'''
from __future__ import print_function

from general_utils import get_random_string, DaosTestError
from general_utils import is_pool_rebuild_complete
from daos_api import DaosApiError

import time


def continuous_io(container, seconds):
    """Perform a combination of reads/writes for the specified time period.

    Args:
        container (DaosContainer): contianer in which to write the data
        seconds (int): how long to write data

    Returns:
        int: number of bytes written to the container

    Raises:
        ValueError: if a data mismatch is detected

    """
    finish_time = time.time() + seconds
    oid = None
    total_written = 0
    size = 500

    while time.time() < finish_time:
        # make some stuff up
        dkey = get_random_string(5)
        akey = get_random_string(5)
        data = get_random_string(size)

        # write it then read it back
        oid, epoch = container.write_an_obj(data, size, dkey, akey, oid, 5)
        data2 = container.read_an_obj(size, dkey, akey, oid, epoch)

        # verify it came back correctly
        if data != data2.value:
            raise ValueError("Data mismatch in ContinousIo")

        # collapse down the commited epochs
        container.consolidate_epochs()

        total_written += size

    return total_written


def write_until_full(container):
    """Write until we get enospace back.

    Args:
        container (DaosContainer): contianer in which to write the data

    Returns:
        int: number of bytes written to the container

    """
    total_written = 0
    size = 2048
    _oid = None

    try:
        while True:
            # make some stuff up and write
            dkey = get_random_string(5)
            akey = get_random_string(5)
            data = get_random_string(size)

            _oid, _epoch = container.write_an_obj(data, size, dkey, akey)
            total_written += size

            # collapse down the commited epochs
            container.slip_epoch()

    except ValueError as exp:
        print(exp)

    return total_written


def write_quantity(container, size_in_bytes):
    """Write a specific number of bytes.

    Note:
        The minimum amount that will be written is 2048 bytes.

    Args:
        container (DaosContainer): which container to write to, it should be in
            an open state prior to the call
        size_in_bytes (int): total number of bytes to be written, although no
            less that 2048 will be written.

    Returns:
        int: number of bytes written to the container

    """
    total_written = 0
    size = 2048
    _oid = None

    try:
        while total_written < size_in_bytes:

            # make some stuff up and write
            dkey = get_random_string(5)
            akey = get_random_string(5)
            data = get_random_string(size)

            _oid, _epoch = container.write_an_obj(data, size, dkey, akey)
            total_written += size

            # collapse down the commited epochs
            container.slip_epoch()

    except ValueError as exp:
        print(exp)

    return total_written


def write_single_objects(
        container, obj_qty, rec_qty, akey_size, dkey_size, data_size, rank,
        object_class, log=None):
    """Write random single objects to the container.

    Args:
        container (DaosContainer): the container in which to write objects
        obj_qty (int): the number of objects to create in the container
        rec_qty (int): the number of records to create in each object
        akey_size (int): the akey length
        dkey_size (int): the dkey length
        data_size (int): the length of data to write in each record
        rank (int): the server rank to which to write the records
        log (DaosLog|None): object for logging messages

    Returns:
        list: a list of dictionaries containing the object, transaction
            number, and data written to the container

    Raises:
        DaosTestError: if an error is detected writing the objects or
            verifying the write of the objects

    """
    if log:
        log.info("Creating objects in the container")
    object_list = []
    for obj in range(obj_qty):
        object_list.append({"obj": None, "txn": None, "record": []})
        for _ in range(rec_qty):
            akey = get_random_string(
                akey_size,
                [record["akey"] for record in object_list[obj]["record"]])
            dkey = get_random_string(
                dkey_size,
                [record["dkey"] for record in object_list[obj]["record"]])
            data = get_random_string(data_size)
            object_list[obj]["record"].append(
                {"akey": akey, "dkey": dkey, "data": data})

            # Write single data to the container
            try:
                (object_list[obj]["obj"], object_list[obj]["txn"]) = \
                    container.write_an_obj(
                        data, len(data), dkey, akey, object_list[obj]["obj"],
                        rank, object_class)
            except DaosApiError as error:
                raise DaosTestError(
                    "Error writing data (dkey={}, akey={}, data={}) to "
                    "the container: {}".format(dkey, akey, data, error))

            # Verify the single data was written to the container
            data_read = read_single_objects(
                container, data_size, dkey, akey, object_list[obj]["obj"],
                object_list[obj]["txn"])
            if data != data_read:
                raise DaosTestError(
                    "Written data confirmation failed:"
                    "\n  wrote: {}\n  read:  {}".format(data, data_read))

    return object_list


def read_single_objects(container, size, dkey, akey, obj, txn):
    """Read data from the container.

    Args:
        container (DaosContainer): the container from which to read objects
        size (int): amount of data to read
        dkey (str): dkey used to access the data
        akey (str): akey used to access the data
        obj (object): object to read
        txn (int): transaction number

    Returns:
        str: data read from the container

    Raises:
        DaosTestError: if an error is dectected reading the objects

    """
    try:
        data = container.read_an_obj(size, dkey, akey, obj, txn)
    except DaosApiError as error:
        raise DaosTestError(
            "Error reading data (dkey={}, akey={}, size={}) from the "
            "container: {}".format(dkey, akey, size, error))
    return data.value


def write_array_objects(
        container, obj_qty, rec_qty, akey_size, dkey_size, data_size, rank,
        object_class, log=None):
    """Write array objects to the container.

    Args:
        container (DaosContainer): the container in which to write objects
        obj_qty (int): the number of objects to create in the container
        rec_qty (int): the number of records to create in each object
        akey_size (int): the akey length
        dkey_size (int): the dkey length
        data_size (int): the length of data to write in each record
        rank (int): the server rank to which to write the records
        log (DaosLog|None): object for logging messages

    Returns:
        list: a list of dictionaries containing the object, transaction
            number, and data written to the container

    Raises:
        DaosTestError: if an error is detected writing the objects or
            verifying the write of the objects

    """
    if log:
        log.info("Creating objects in the container")
    object_list = []
    for obj in range(obj_qty):
        object_list.append({"obj": None, "txn": None, "record": []})
        for _ in range(rec_qty):
            akey = get_random_string(
                akey_size,
                [record["akey"] for record in object_list[obj]["record"]])
            dkey = get_random_string(
                dkey_size,
                [record["dkey"] for record in object_list[obj]["record"]])
            data = [get_random_string(data_size) for _ in range(data_size)]
            object_list[obj]["record"].append(
                {"akey": akey, "dkey": dkey, "data": data})

            # Write the data to the container
            try:
                object_list[obj]["obj"], object_list[obj]["txn"] = \
                    container.write_an_array_value(
                        data, dkey, akey, object_list[obj]["obj"], rank,
                        object_class)
            except DaosApiError as error:
                raise DaosTestError(
                    "Error writing data (dkey={}, akey={}, data={}) to "
                    "the container: {}".format(dkey, akey, data, error))

            # Verify the data was written to the container
            data_read = read_array_objects(
                container, data_size, data_size + 1, dkey, akey,
                object_list[obj]["obj"], object_list[obj]["txn"])
            if data != data_read:
                raise DaosTestError(
                    "Written data confirmation failed:"
                    "\n  wrote: {}\n  read:  {}".format(data, data_read))

    return object_list


def read_array_objects(container, size, items, dkey, akey, obj, txn):
    """Read data from the container.

    Args:
        container (DaosContainer): the container from which to read objects
        size (int): number of arrays to read
        items (int): number of items in each array to read
        dkey (str): dkey used to access the data
        akey (str): akey used to access the data
        obj (object): object to read
        txn (int): transaction number

    Returns:
        str: data read from the container

    Raises:
        DaosTestError: if an error is dectected reading the objects

    """
    try:
        data = container.read_an_array(
            size, items, dkey, akey, obj, txn)
    except DaosApiError as error:
        raise DaosTestError(
            "Error reading data (dkey={}, akey={}, size={}, items={}) "
            "from the container: {}".format(
                dkey, akey, size, items, error))
    return [item[:-1] for item in data]


def read_during_rebuild(
        container, pool, log, written_objects, data_size, read_method):
    """Read all the written data while rebuild is still in progress.

    Args:
        container (DaosContainer): the container from which to read objects
        pool (DaosPool): pool for which to determine if rebuild is complete
        log (logging): logging object used to report the pool status
        written_objects (list): record data type to write/read
        data_size (int):
        read_method (object): function to call to read the data, e.g.
            read_single_objects or read_array_objects

    Returns:
        None

    Raises:
        DaosTestError: if the rebuild completes before the read is complete
            or read errors are detected

    """
    index_x = 0
    index_y = 0
    incomplete = True
    failed_reads = False
    while not is_pool_rebuild_complete(pool, log) and incomplete:
        incomplete = index_x < len(written_objects)
        if incomplete:
            # Read the data from the previously written record
            record_info = {
                "container": container,
                "size": data_size,
                "dkey": written_objects[index_x]["record"][index_y]["dkey"],
                "akey": written_objects[index_x]["record"][index_y]["akey"],
                "obj": written_objects[index_x]["obj"],
                "txn": written_objects[index_x]["txn"]}
            if read_method.__name__ == "read_array_objects":
                record_info["items"] = data_size + 1
            read_data = read_method(**record_info)

            # Verify the data just read matches the original data written
            record_info["data"] = \
                written_objects[index_x]["record"][index_y]["data"]
            if record_info["data"] != read_data:
                failed_reads = True
                log.error(
                    "<obj: %s, rec: %s>: Failed reading data "
                    "(dkey=%s, akey=%s):\n  read:  %s\n  wrote: %s",
                    index_x, index_y, record_info["dkey"], record_info["akey"],
                    read_data, record_info["data"])
            else:
                log.info(
                    "<obj: %s, rec: %s>: Passed reading data "
                    "(dkey=%s, akey=%s)",
                    index_x, index_y, record_info["dkey"], record_info["akey"])

            # Read the next record in this object or the next object
            index_y += 1
            if index_y >= len(written_objects[index_x]["record"]):
                index_x += 1
                index_y = 0

    # Verify that all of the objects and records were read successfully
    # while the rebuild was still active
    if incomplete:
        raise DaosTestError(
            "Rebuild completed before all the written data could be read")
    elif failed_reads:
        raise DaosTestError("Errors detected reading data during rebuild")