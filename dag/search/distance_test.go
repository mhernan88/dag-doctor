package search

// import (
//     "testing"
//     "github.com/sirupsen/logrus"
// )

// func TestDistanceCalculator__DistanceToStart(t *testing.T) {
//     l := logrus.New()
//     l.SetLevel(logrus.TraceLevel)
//     dc := NewDistanceCalculator(l)
//
//     nodes := setupDAG()
//     endPtr := nodes["node5a"]
//
//     result, err := dc.DistanceToStart(endPtr, 99)
//     if err != nil {
//         t.Error(err)
//     }
//     if result != 4 {
//         t.Errorf("want=4, got=%d", result)
//     }
// }
//
// func TestDistanceCalculator__DistanceToEnd(t *testing.T) {
//     l := logrus.New()
//     l.SetLevel(logrus.TraceLevel)
//     dc := NewDistanceCalculator(l)
//
//     nodes := setupDAG()
//     startPtr := nodes["node1"]
//
//     result, err := dc.DistanceToEnd(startPtr, 99)
//     if err != nil {
//         t.Error(err)
//     }
//     if result != 4 {
//         t.Errorf("want=4, got=%d", result)
//     }
// }
//
// func TestDistanceCalculator__Midpoints(t *testing.T) {
//     l := logrus.New()
//     l.SetLevel(logrus.TraceLevel)
//     dc := NewDistanceCalculator(l)
//
//     nodes := setupDAG()
//     startPtr := nodes["node1"]
//
//     midpoints, err := dc.Midpoints([]*Node{startPtr}, 99)
//     if err != nil {
//         t.Error(err)
//         return
//     }
//     if len(nodes) == 0 {
//         t.Error("nodes input was modified")
//         return
//     }
//     if len(midpoints) == 0 {
//         t.Error("dc found no midpoints")
//         return
//     }
//     if len(midpoints) != 1 {
//         t.Errorf("dc found wrong number of midpoints; got=%d, want=1", len(midpoints))
//         return
//     }
//     if midpoints[0] != "node3a" {
//         t.Errorf("dc found wrong midpoint; got='%s', want='node3a'", midpoints[0])
//         return
//     }
// }
