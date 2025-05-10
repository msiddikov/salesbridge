package meta

func init() {

}

func getTestClient() Client {
	testToken := "EAAK2EsTlzmEBO2anEplLwLyVTu0FCB3pCjgqxYZAc6HIxZAEco5zmRyhN15DHaqKU6IrXQnkljGyIOayRAKYSFZCjfFDiDp7gVBrBm7vbCHLJz43kupv1KhfN1ncDwmPHSqHVWxdyo9FbQaHP8KurPBrS3j0r6aKHKvhQVhNRKR0J3B7DIqoPGM"
	s := Service{}
	cli, _ := s.NewClient(testToken)
	return cli
}
