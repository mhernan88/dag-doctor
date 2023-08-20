package io

import (
    "slices"
    "testing"
    "github.com/sirupsen/logrus"
)

func TestLoadNodes(t *testing.T) {
    l := logrus.New()
    l.SetLevel(logrus.TraceLevel)

    roots, err := LoadNodes("../../dag.json", l)
    if err != nil {
        t.Error(err)
    }

    if len(roots) != 3 {
        t.Errorf("expected 3 roots, got %d", len(roots))
    }

    var names []string
    for _, root := range roots {
        t.Logf("root %s found", root.Name)
        names = append(names, root.Name)
    }

    expecteds := []string{
        "preprocess_companies_and_employees",
        "preprocess_shuttles_and_routes",
        "preprocess_reviews_and_ratings",
    }

    for _, expected := range expecteds {
        t.Logf("reviewing expected: %s", expected)
        if !slices.Contains(names, expected) {
            t.Errorf("roots missing %s", expected)
        }
    }
}
