/*
 * (C) Copyright 2021 Intel Corporation.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 * http://www.apache.org/licenses/LICENSE-2.0
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

/*
 * This utility shows metrics from the specified IO Server
 */

#include <getopt.h>
#include <string.h>
#include "gurt/telemetry_common.h"
#include "gurt/telemetry_consumer.h"

void
print_my_children(uint64_t *shmem_root, struct d_tm_node_t *node, int filter,
		  int level, FILE *stream)
{
	if ((node == NULL) || (stream == NULL))
		return;

	if (node->dtn_type & filter)
		d_tm_print_node(shmem_root, node, level, stream);

	node = node->dtn_child;
	node = d_tm_conv_ptr(shmem_root, node);
	if (node == NULL)
		return;

	print_my_children(shmem_root, node, filter, level + 1, stream);
	node = node->dtn_sibling;
	node = d_tm_conv_ptr(shmem_root, node);
	while (node != NULL) {
		print_my_children(shmem_root, node, filter, level + 1, stream);
		node = node->dtn_sibling;
		node = d_tm_conv_ptr(shmem_root, node);
	}
}


static void
print_usage(const char *prog_name)
{
	printf("Usage: %s [optional arguments]\n"
	       "\n"
	       "--rank, r\n"
	       "\tShow telemetry data from this IO Server rank (default 0)\n"
	       "--path, p\n"
	       "\tDisplay metrics at or below the specified path\n"
	       "\tDefault is root directory\n"
	       "--iterations, i\n"
	       "\tSpecifies the number of iterations to show "
	       "(default continuous)\n"
	       "--delay, D\n"
	       "\tDelay in seconds between each iteration\n"
	       "\tDefault is 1 second\n\n"
	       "Filter the displayed data by specifying one or more filters:\n"
	       "\tDefault is all filters\n\n"
	       "--counter, -c\n"
	       "\tInclude counters\n"
	       "--duration, -d\n"
	       "\tInclude durations\n"
	       "--timestamp, -t\n"
	       "\tInclude timestamps\n"
	       "--snapshot, -s\n"
	       "\tInclude timer snapshots\n"
	       "--gauge, -g\n"
	       "\tInclude gauges\n",
	       prog_name);
}

int
main(int argc, char **argv)
{
	struct d_tm_node_t	*root = NULL;
	struct d_tm_node_t	*node = NULL;
	uint64_t		*shmem_root = NULL;
	char			dirname[D_TM_MAX_NAME_LEN] = {0};
	int			rank = 0;
	int			iteration = 0;
	int			num_iter = 0;
	int			filter = 0;
	int			delay = 1;
	int			c;

	sprintf(dirname, "/");

	/********************* Parse user arguments *********************/
	while (1) {
		static struct option long_options[] = {
			{"rank", required_argument, NULL, 'r'},
			{"counter", no_argument, NULL, 'c'},
			{"duration", no_argument, NULL, 'd'},
			{"timestamp", no_argument, NULL, 't'},
			{"snapshot", no_argument, NULL, 's'},
			{"gauge", no_argument, NULL, 'g'},
			{"iterations", required_argument, NULL, 'i'},
			{"path", required_argument, NULL, 'p'},
			{"delay", required_argument, NULL, 'D'},
			{"help", no_argument, NULL, 'h'},
			{NULL, 0, NULL, 0}
		};

		c = getopt_long_only(argc, argv, "r:cdtsgi:p:D:h",
				     long_options, NULL);
		if (c == -1)
			break;

		switch (c) {
		case 'r':
			rank = atoi(optarg);
			break;
		case 'c':
			filter |= D_TM_COUNTER;
			break;
		case 'd':
			filter |= D_TM_DURATION;
			break;
		case 't':
			filter |= D_TM_TIMESTAMP;
			break;
		case 's':
			filter |= D_TM_TIMER_SNAPSHOT;
			break;
		case 'g':
			filter |= D_TM_GAUGE;
			break;
		case 'i':
			num_iter = atoi(optarg);
			break;
		case 'p':
			strncpy(dirname, optarg, D_TM_MAX_NAME_LEN);
			break;
		case 'D':
			delay = atoi(optarg);
			break;
		case 'h':
		case '?':
		default:
			print_usage(argv[0]);
			exit(0);
		}
	}

	if (filter == 0)
		filter = D_TM_COUNTER | D_TM_DURATION | D_TM_TIMESTAMP |
			 D_TM_TIMER_SNAPSHOT | D_TM_GAUGE;

	shmem_root = d_tm_get_shared_memory(rank);
	if (!shmem_root)
		goto failure;

	root = d_tm_get_root(shmem_root);
	if (!root)
		goto failure;

	node = d_tm_find_metric(shmem_root, dirname);
	if (node != NULL)
		root = node;

	while ((num_iter == 0) || (iteration < num_iter)) {
		printf("\nIteration: %d\n", iteration);
		printf("There are %" PRIu64 " matching metrics underneath the"
		       " directory %s\n",
		       d_tm_count_metrics(shmem_root, root, filter),
		       dirname[0] != '\0' ? dirname : "/");

		print_my_children(shmem_root, root, filter | D_TM_DIRECTORY, 0,
				  stdout);
		iteration++;
		sleep(delay);
			printf("\n\n");
	}

	return 0;

failure:
	printf("Unable to attach to the shared memory for the rank: %d\n"
	       "Make sure to run the IO server with the same rank to initialize"
	       " the shared memory and populate it with metrics.\n",
	       rank);
	return -1;
}
