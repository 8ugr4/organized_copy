package pkg

import (
	"github.com/stretchr/testify/require"
	"testing"
)

func Test_initExifTool(t *testing.T) {
	o := GetNewOperator()
	require.Nil(t, o.Storage.Exif)

	err := o.initExifTool()
	require.NoError(t, err)
	require.NotNil(t, o.Storage.Exif)
}

//func Test_getFileDate(t *testing.T) {
//	o := GetNewOperator()
//	fp := "../testDir/dummy2.jpg"
//	regexPattern := "*"
//	periodType := "year"
//	res, err := o.getFileDate(fp, regexPattern, periodType)
//	fmt.Println(res, err)
//}
