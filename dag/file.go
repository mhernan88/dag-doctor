package dag
//
// import (
//     "encoding/csv"
//     "os"
//     "io"
//     "github.com/fatih/color"
//     "github.com/rodaine/table"
// )
//
//
// func (pf PipelineFile) InitTable(names []string) table.Table {
//     headerFmt := color.New(color.FgGreen, color.Underline).SprintfFunc()
//     columnFmt := color.New(color.FgYellow).SprintfFunc()
//
//     interfaceNames := make([]interface{}, len(names))
//     for i, v := range names {
//         interfaceNames[i] = v
//     }
//
//     tbl := table.New(interfaceNames...)
//     tbl.WithHeaderFormatter(headerFmt).WithFirstColumnFormatter(columnFmt)
//
//     return tbl
// }
//
// func (pf PipelineFile) DisplayTable(tbl table.Table, rows [][]string) {
//     var interfaceRow []interface{}
//     for _, row := range rows {
//         interfaceRow = make([]interface{}, len(row))
//         for i, v := range row {
//             interfaceRow[i] = v
//         }
//         tbl.AddRow(interfaceRow...)
//     }
//     tbl.Print()
// }
//
// func (pf PipelineFile) LoadCSV(n int) ([]string, [][]string, error) {
//     f, err := os.Open(pf.Filename)
//     if err != nil {
//         return nil, nil, err
//     }
//     defer f.Close()
//
//     r := csv.NewReader(f)
//
//     header, err := r.Read()
//     if err != nil {
//         return nil, nil, err
//     }
//
//     var rows [][]string
//
//     i := 0
//     for {
//         row, err := r.Read()
//         if err == io.EOF {
//             break
//         }
//         rows = append(rows, row)
//
//         i++
//         if i >= n {
//             break
//         }
//     }
//     return header, rows, nil
// }
//
// func (pf PipelineFile) LoadAndDisplay(n int) error {
//     header, rows, err := pf.LoadCSV(n)
//     if err != nil {
//         return err
//     }
//
//     tbl := pf.InitTable(header)
//     pf.DisplayTable(tbl, rows)
//     return nil
// }
