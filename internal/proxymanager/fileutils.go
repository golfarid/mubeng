package proxymanager

import (
	"bufio"
	"bytes"
	"ktbs.dev/mubeng/pkg/helper"
	"os"
)

func (p *ProxyManager) ReadProxies() ([]string, error) {
	file, err := os.Open(p.filepath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var proxies []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		proxyUrl := helper.Eval(scanner.Text())
		proxies = append(proxies, proxyUrl)
	}

	return proxies, nil
}

func (p *ProxyManager) WriteProxies(urls []string) error {
	file, err := os.OpenFile(p.filepath, os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer file.Close()

	var bs []byte
	buf := bytes.NewBuffer(bs)

	for i := 0; i < len(urls); i++ {
		_, err := buf.WriteString(urls[i] + "\n")
		if err != nil {
			return err
		}
	}

	_, err = file.Write(buf.Bytes())
	if err != nil {
		return err
	}

	return nil
}

func (p *ProxyManager) DeleteProxy(index int) error {
	file, err := os.Open(p.filepath)
	if err != nil {
		return err
	}

	var bs []byte
	buf := bytes.NewBuffer(bs)

	currentIndex := 0
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		if currentIndex != index {
			_, err := buf.WriteString(line + "\n")
			if err != nil {
				return err
			}
		}

		currentIndex += 1
	}
	if err := scanner.Err(); err != nil {
		return err
	}
	err = file.Close()
	if err != nil {
		return err
	}

	file, err = os.Create(p.filepath)
	if err != nil {
		return err
	}

	_, err = buf.WriteTo(file)
	if err != nil {
		return err
	}

	err = file.Close()
	if err != nil {
		return err
	}
	return nil
}
