// Copyright 2021 dfuse Platform Inc.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package types

import (
	"fmt"
	"io"
	"time"

	"github.com/streamingfast/bstream"
)

const (
	dbinContentType = "ARE" // Arweave
)

func init() {
	bstream.GetBlockWriterFactory = bstream.BlockWriterFactoryFunc(blockWriterFactory)
	bstream.GetBlockReaderFactory = bstream.BlockReaderFactoryFunc(blockReaderFactory)
	bstream.GetBlockDecoder = bstream.BlockDecoderFunc(BlockDecoder)
	bstream.GetBlockWriterHeaderLen = 10
	bstream.GetBlockPayloadSetter = bstream.MemoryBlockPayloadSetter
	bstream.GetMemoizeMaxAge = 200 * 15 * time.Second

	// We want to panic in here to enforce validation in any component that uses this package,
	// instead of running validation in multiple places.
	if err := bstream.ValidateRegistry(); err != nil {
		panic(err)
	}
}

func blockReaderFactory(reader io.Reader) (bstream.BlockReader, error) {
	return bstream.NewDBinBlockReader(reader, func(contentType string, version int32) error {
		if contentType != dbinContentType && version != 1 {
			return fmt.Errorf("reader only knows about %s block kind at version 1, got %s at version %d", dbinContentType, contentType, version)
		}

		return nil
	})
}

func blockWriterFactory(writer io.Writer) (bstream.BlockWriter, error) {
	return bstream.NewDBinBlockWriter(writer, dbinContentType, 1)
}
