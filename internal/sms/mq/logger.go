// Copyright 2022 Lingfei Kong <colin404@foxmail.com>. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file. The original repo for
// this file is https://github.com/Rosas99/smsx.
//

package mq

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/segmentio/kafka-go"

	"github.com/Rosas99/smsx/pkg/log"
	genericoptions "github.com/Rosas99/smsx/pkg/options"
)

// kafkaLogger is a log.Logger implementation that writes log messages to Kafka.
type KafkaLogger struct {
	// enabled is an atomic boolean indicating whether the logger is enabled.
	enabled int32
	// writer is the Kafka writer used to write log messages.
	writer *kafka.Writer
}

// todo 这里改成短信发送历史
// AuditMessage is the message structure for log messages.
type AuditMessage struct {
	Matcher   string     `protobuf:"bytes,1,opt,name=matcher,proto3" json:"matcher,omitempty"`
	Request   []any      `protobuf:"bytes,2,opt,name=request,proto3" json:"request,omitempty"`
	Result    bool       `protobuf:"bytes,3,opt,name=result,proto3" json:"result,omitempty"`
	Explains  [][]string `protobuf:"bytes,4,opt,name=explains,proto3" json:"explains,omitempty"`
	Timestamp int64      `protobuf:"bytes,5,opt,name=timestamp,proto3" json:"timestamp,omitempty"`
}

// NewLogger creates a new kafkaLogger instance.
func NewLogger(kafkaOpts *genericoptions.KafkaOptions) (*KafkaLogger, error) {
	writer, err := kafkaOpts.Writer()
	if err != nil {
		return nil, err
	}

	return &KafkaLogger{writer: writer}, nil
}

// LogModel writes a log message for the policy model.
func (l *KafkaLogger) LogHistory(test string) string {

	//message := AuditMessage{
	//	Timestamp: time.Now().Unix(),
	//}

	// 这里先转成json
	out, _ := json.Marshal(test)
	// 因为writer传的都是二进制
	// 这里json转成二进制数组
	fmt.Println(test)
	if err := l.writer.WriteMessages(context.Background(), kafka.Message{Value: out}); err != nil {
		log.Errorw(err, "Failed to write kafka messages")
	} else {
		fmt.Println(string(out))
	}
	return string(out)
}

// log others
