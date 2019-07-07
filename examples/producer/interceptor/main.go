/*
Licensed to the Apache Software Foundation (ASF) under one or more
contributor license agreements.  See the NOTICE file distributed with
this work for additional information regarding copyright ownership.
The ASF licenses this file to You under the Apache License, Version 2.0
(the "License"); you may not use this file except in compliance with
the License.  You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

// Package main implements a producer with user custom interceptor.
package main

import (
	"context"
	"fmt"
	"os"
	"strconv"

	"github.com/apache/rocketmq-client-go/internal/producer"
	"github.com/apache/rocketmq-client-go/primitive"
)

func main() {
	nameServerAddr := "127.0.0.1:9876"
	p, _ := producer.NewProducer(nameServerAddr, primitive.WithRetry(2),
		primitive.WithChainProducerInterceptor(UserFirstInterceptor(), UserSecondInterceptor()))
	err := p.Start()
	if err != nil {
		fmt.Printf("start producer error: %s", err.Error())
		os.Exit(1)
	}
	for i := 0; i < 10; i++ {
		res, err := p.SendSync(context.Background(), &primitive.Message{
			//Topic: "test",
			Topic: "TopicTest",
			Body:  []byte("Hello RocketMQ Go Client!"),
			Properties: map[string]string{"order": strconv.Itoa(i)},
		})

		if err != nil {
			fmt.Printf("send message error: %s\n", err)
		} else {
			fmt.Printf("send message success: result=%s\n", res.String())
		}
	}
	err = p.Shutdown()
	if err != nil {
		fmt.Printf("shundown producer error: %s", err.Error())
	}
}

func UserFirstInterceptor() primitive.PInterceptor {
	return func(ctx context.Context, req, reply interface{}, next primitive.PInvoker) error {
		fmt.Printf("user first interceptor before invoke: req:%v\n", req)
		err := next(ctx, req, reply)
		fmt.Printf("user first interceptor after invoke: req: %v, reply: %v \n", req, reply)
		return err
	}
}

func UserSecondInterceptor() primitive.PInterceptor {
	return func(ctx context.Context, req, reply interface{}, next primitive.PInvoker) error {
		fmt.Printf("user second interceptor before invoke: req: %v\n", req)
		err := next(ctx, req, reply)
		fmt.Printf("user second interceptor after invoke: req: %v, reply: %v \n", req, reply)
		return err
	}
}
