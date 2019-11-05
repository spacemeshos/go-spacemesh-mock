package main

import (
	"context"
	"fmt"
	"strconv"
	"testing"
	"time"

	"github.com/spacemeshos/go-spacemesh-mock/api/nmpb"
	"github.com/spacemeshos/go-spacemesh-mock/integration"

	"github.com/spacemeshos/smutil/log"
	"github.com/stretchr/testify/require"
)

// harnessTestCase utilizes an instance of the Harness to
// exercise functionality.
type harnessTestCase struct {
	name string
	test func(h *integration.Harness, assert *require.Assertions, ctx context.Context)
}

var testCases = []*harnessTestCase {
	{name: "Mock test", test: testMock},
}

// NewHarness creates and initializes a new instance of Harness.
func newHarness(req *require.Assertions, cfg *integration.ServerConfig) * integration.Harness {
	h, err := integration.NewHarness(cfg)
	req.NoError(err)
	req.NotNil(h)

	go func() {
		for {
			err, more := <-h.ProcessErrors()
			if !more {
				return
			}
			req.Fail("mockNode server has finished with errors", err)
		}
	}()

	return h
}

func TestHarness(t *testing.T) {
	assert := require.New(t)

	srcCodePath := "."
	cfg, err := integration.DefaultConfig(srcCodePath)
	assert.NoError(err)

	h := newHarness(assert, cfg)

	defer func() {
		err := h.TearDown()
		assert.NoError(err, "failed to teardown mockNode's harness")
		t.Logf("harness teared down")
	}()

	for _, testCase := range testCases {
		success := t.Run(testCase.name, func(t1 *testing.T) {
			fmt.Println("running:", testCase.name + "\n")
			ctx, _ := context.WithTimeout(context.Background(), 1*time.Second)
			testCase.test(h, assert, ctx)
		})

		if !success {
			break
		}
	}
}

func testMock(h *integration.Harness, assert *require.Assertions, ctx context.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	numOfRounds := 10
	proof := "This is the proof for round %v broadcast"

	for i := 1; i <= numOfRounds; i++ {
		res, err := broadcast(h, ctx, []byte(fmt.Sprintf(proof, i)))
		assert.Nil(err)
		fmt.Printf(proof + "\n", i)
		fmt.Printf("broadcast validation msg for round %v: %v\n", i, res)
	}

	for r := 1; r <= numOfRounds; r++ {
		content, err := getProof(h, ctx, r)
		assert.Nil(err)
		fmt.Println("round", r, "returned content:", string(content))
	}
}

// broadcast poets' proof and returns result msg type of a string and error.
func broadcast(h *integration.Harness, ctx context.Context, proof []byte) (string, error) {
	msg, err := h.BroadcastPoet(ctx, &nmpb.BinaryMessage{Data: proof})
	if err != nil {
		log.Error("could not broadcast proof:", err)
		return "", err
	}

	return msg.Value, nil
}

// get proof by round index, return byte stream of the proof and error.
func getProof(h *integration.Harness, ctx context.Context, roundNum int) ([]byte, error) {
	binMsg, err := h.GetProof(ctx, &nmpb.SimpleMessage{Value: strconv.Itoa(roundNum)})
	if err != nil {
		log.Error("could not retrieve data: %v", err)
		return nil, err
	}

	return binMsg.Data, nil
}
