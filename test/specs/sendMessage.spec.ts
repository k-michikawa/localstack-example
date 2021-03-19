import * as AWS from "aws-sdk";
import { CredentialsOptions } from "aws-sdk/lib/credentials";
import { ListQueuesRequest, ReceiveMessageRequest } from "aws-sdk/clients/sqs";
import api from "../helpers/api";

AWS.config.update({
  region: "ap-northeast-1",
});

const endpoint = "http://localhost:4566";

const QUEUE_URL = [endpoint, "localstack", "example"].join("/");

const credentials: CredentialsOptions = {
  accessKeyId: "dummy",
  secretAccessKey: "dummy",
};

const sqs = new AWS.SQS({
  credentials,
  apiVersion: "2012-11-05",
  endpoint: new AWS.Endpoint(endpoint),
});

describe("send-message endpoint test", () => {
  test("success", async () => {
    const listQueueParams: ListQueuesRequest = {
      MaxResults: 1,
    };
    const queues = await sqs.listQueues(listQueueParams).promise();
    console.log(queues);

    // check queue exists
    if (!queues.QueueUrls || queues.QueueUrls.length < 1) {
      throw new Error("Queue not found");
    }
    console.log({ queueUrl: queues.QueueUrls[0] });

    const message = "test";
    const apiRes = await api.post("/send-message", { message });

    if (apiRes.status !== 200) {
      throw new Error("Faild call api");
    }

    const receiveMessageParams: ReceiveMessageRequest = {
      QueueUrl: QUEUE_URL,
      AttributeNames: ["All"],
      MaxNumberOfMessages: 10,
      VisibilityTimeout: 60,
      WaitTimeSeconds: 5,
    };
    const receiveMessageRes = await sqs
      .receiveMessage(receiveMessageParams)
      .promise();

    if (!receiveMessageRes.Messages) {
      throw new Error("Message not found");
    }

    expect(receiveMessageRes.Messages.length).toEqual(1);

    const enqueuedMessage = receiveMessageRes.Messages[0];

    expect(enqueuedMessage!.MessageId).toEqual(apiRes.data);
    expect(enqueuedMessage!.Body).toEqual(message);

    if (!enqueuedMessage.ReceiptHandle) {
      throw new Error("ReceiptHandle is undefined");
    }
    // 対象のメッセージを消しておく
    await sqs
      .deleteMessage({
        QueueUrl: QUEUE_URL,
        ReceiptHandle: enqueuedMessage.ReceiptHandle,
      })
      .promise();
  });
});
