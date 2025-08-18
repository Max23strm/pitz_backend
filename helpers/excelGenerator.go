package helpers

import (
	"fmt"
	"strconv"

	"github.com/Max23strm/pitz-backend/models"
	"github.com/xuri/excelize/v2"
)

func CreatePaymentExcel(payments models.PaymentFileRows, monthlyPayments models.MonthlyFileRows) (*excelize.File, error) {

	f := excelize.NewFile()
	defer func() {
		if err := f.Close(); err != nil {
			fmt.Println(err)
		}
	}()

	// Cabecera de columnas
	headerStyle, err := f.NewStyle(&excelize.Style{
		Fill: excelize.Fill{Type: "pattern", Color: []string{"0C5C7A"}, Pattern: 1},
		Font: &excelize.Font{
			Bold:  true,
			Color: "FFFFFF",
		},
		NumFmt: 4,
	})
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	numberStyle, err := f.NewStyle(&excelize.Style{
		NumFmt: 4,
	})
	dateStyle, err := f.NewStyle(&excelize.Style{
		NumFmt: 15,
	})
	// monthStyle, err := f.NewStyle(&excelize.Style{
	// 	NumFmt: 18,
	// })
	if err != nil {
		fmt.Println("error formateando numero" + err.Error())
		return nil, err
	}

	f.SetCellValue("Sheet1", "A1", "Nombres")
	f.SetCellValue("Sheet1", "B1", "Apellidos")
	f.SetCellValue("Sheet1", "C1", "Fecha de pago")
	f.SetCellValue("Sheet1", "D1", "Tipo de pago")
	f.SetCellValue("Sheet1", "E1", "Monto")
	f.SetCellValue("Sheet1", "F1", "Comentario")
	f.SetCellValue("Sheet1", "I1", "Mes")
	f.SetCellValue("Sheet1", "J1", "Total de Mes")
	f.SetCellStyle("Sheet1", "A1", "F1", headerStyle)
	f.SetCellStyle("Sheet1", "I1", "J1", headerStyle)

	totalAmount := 0.00

	totalRow := strconv.Itoa(len(payments) + 3)
	for i, payment := range payments {
		currentRow := strconv.Itoa(i + 2)

		f.SetCellValue("Sheet1", "A"+currentRow, payment.Player_name)
		f.SetCellValue("Sheet1", "B"+currentRow, payment.Player_last_name)
		f.SetCellValue("Sheet1", "C"+currentRow, payment.Payment_date)
		f.SetCellValue("Sheet1", "D"+currentRow, payment.Payment_name)
		f.SetCellValue("Sheet1", "E"+currentRow, payment.Amount)
		f.SetCellValue("Sheet1", "F"+currentRow, payment.Comment)

		f.SetCellStyle("Sheet1", "E"+currentRow, "E"+currentRow, numberStyle)
		f.SetCellStyle("Sheet1", "C"+currentRow, "C"+currentRow, dateStyle)
		totalAmount = totalAmount + payment.Amount

	}

	for i, month := range monthlyPayments {
		currentRow := strconv.Itoa(i + 2)

		f.SetCellValue("Sheet1", "I"+currentRow, month.Month)
		f.SetCellValue("Sheet1", "J"+currentRow, month.Amount)

		f.SetCellStyle("Sheet1", "J"+currentRow, "J"+currentRow, numberStyle)
		// f.SetCellStyle("Sheet1", "I"+currentRow, "I"+currentRow, monthStyle)

	}

	f.SetCellValue("Sheet1", "A"+totalRow, "Total")
	f.SetCellValue("Sheet1", "E"+totalRow, totalAmount)
	f.SetCellStyle("Sheet1", "A"+totalRow, "F"+totalRow, headerStyle)

	return f, nil
	// if err := f.SaveAs("Book1.xlsx"); err != nil {
	// 	fmt.Println(err)
	// }

}
