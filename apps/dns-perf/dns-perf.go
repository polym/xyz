package main

import (
	"context"
	"flag"
	"fmt"
	"net"
	"sync"
	"time"
)

var (
	resolver *net.Resolver
	lock     sync.Mutex
)

func dnsLookUp(domain string, timeout int) error {
	dur := time.Duration(timeout) * time.Second
	ctx, cancel := context.WithTimeout(context.Background(), dur)
	defer cancel()
	_, err := resolver.LookupIPAddr(ctx, domain)
	if err != nil {
		fmt.Println(err)
	}
	return err
}

func main() {
	var (
		domain, host              string
		num, conns, port, timeout int
	)
	flag.StringVar(&domain, "d", "v403.upyun.local", "domain to resolv")
	flag.StringVar(&host, "h", "", "dns server host")
	flag.IntVar(&port, "p", 53, "dns server port")
	flag.IntVar(&num, "n", 1000, "number of requests")
	flag.IntVar(&conns, "c", 50, "connections to request")
	flag.IntVar(&timeout, "timeout", 3, "resolv timeout")
	flag.Parse()

	fmt.Println("============================================")
	if host != "" {
		fmt.Printf("dns server             : %s:%d\n", host, port)
	} else {
		fmt.Printf("dns server             : /etc/resolv.conf\n")
	}
	fmt.Printf("domain to resolv       : %s\n", domain)
	fmt.Printf("number of requests     : %d\n", num)
	fmt.Printf("connections to request : %d\n", conns)
	fmt.Printf("resolv timeout         : %ds\n", timeout)
	fmt.Println("============================================\n")

	if host != "" {
		resolver = &net.Resolver{
			PreferGo: true,
			Dial: func(ctx context.Context, network, address string) (net.Conn, error) {
				d := net.Dialer{}
				return d.DialContext(ctx, "udp", net.JoinHostPort(host, fmt.Sprint(port)))
			},
		}
	} else {
		resolver = net.DefaultResolver
	}

	wg := sync.WaitGroup{}
	qchan := make(chan struct{}, conns*10)
	totalCount, errCount := 0, 0
	for i := 0; i < conns; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for _ = range qchan {
				err := dnsLookUp(domain, timeout)
				lock.Lock()
				totalCount++
				if err != nil {
					errCount++
				}
				lock.Unlock()
			}
		}()
	}

	go func() {
		timeBegin, timeBefore, totalBefore := time.Now(), time.Now(), 0
		time.Sleep(time.Second)
		for {
			lock.Lock()
			tmpTotal := totalCount
			lock.Unlock()
			now := time.Now()
			qps := (tmpTotal - totalBefore) / int(now.Sub(timeBefore)/time.Second)
			qpsTotal := tmpTotal / int(now.Sub(timeBegin)/time.Second)
			fmt.Printf("total: %10d err: %10d qps: %10d total_qps: %10d\n", tmpTotal, errCount, qps, qpsTotal)
			time.Sleep(time.Second)
			timeBefore, totalBefore = now, tmpTotal
		}
	}()

	for i := 0; i < num; i++ {
		qchan <- struct{}{}
	}
	close(qchan)

	wg.Wait()

	fmt.Printf("total: %d err: %d\n", totalCount, errCount)
}
