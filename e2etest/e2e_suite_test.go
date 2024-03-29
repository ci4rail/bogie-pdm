package e2e

import (
	"context"
	"log"
	"math/rand"
	"testing"
	"time"

	"github.com/edgefarm/edgefarm.core/test/pkg/framework"
	"github.com/onsi/ginkgo/v2"
	"github.com/onsi/gomega"

	// register tests
	_ "github.com/ci4rail/bogie-pdm/e2etest/powercycle"
)

func TestKube(t *testing.T) {
	rand.Seed(time.Now().UTC().UnixNano())
	gomega.RegisterFailHandler(ginkgo.Fail)
	err := framework.CreateFramework(context.Background(), nil)
	if err != nil {
		log.Fatalf("Failed to create framework: %v", err)
	}
	ginkgo.RunSpecs(t, "bogie Suite")
}
