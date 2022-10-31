package service

import (
	"github.com/eleven-sh/eleven/entities"
)

func prefixClusterResource(
	clusterNameSlug string,
) func(string) string {

	return func(resourceNameSlug string) string {
		if clusterNameSlug == entities.DefaultClusterName {
			return "eleven-" + resourceNameSlug
		}

		return "eleven-" + clusterNameSlug + "-" + resourceNameSlug
	}
}

func prefixEnvResource(
	clusterNameSlug string,
	envNameSlug string,
) func(string) string {

	return func(resourceNameSlug string) string {
		if clusterNameSlug == entities.DefaultClusterName {
			return "eleven-" + envNameSlug + "-" + resourceNameSlug
		}

		return "eleven-" + clusterNameSlug + "-" + envNameSlug + "-" + resourceNameSlug
	}
}
