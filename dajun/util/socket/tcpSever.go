package main

func server() error {
	endpoint := NewEndpoint()

	endpoint.AddHandleFunc("STRING", handleStrings)
	endpoint.AddHandleFunc("GOB", handleGob)

	return endpoint.Listen()
}
