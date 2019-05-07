#!/bin/bash
protoc --api_linter_out=out_format=yaml,cfg_file=test_cfg.json:. test.proto