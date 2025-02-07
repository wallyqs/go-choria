package cmd

import (
	"fmt"
	"sync"
	"time"

	"github.com/choria-io/go-choria/choria"
	"github.com/choria-io/go-choria/srvcache"
)

type tSubCommand struct {
	command
	topic string
	raw   bool
}

func (s *tSubCommand) Setup() (err error) {
	if tool, ok := cmdWithFullCommand("tool"); ok {
		s.cmd = tool.Cmd().Command("sub", "Subscribe to middleware topics")
		s.cmd.Arg("topic", "The topic to subscribe to").Required().StringVar(&s.topic)
		s.cmd.Flag("raw", "Display raw messages one per line without timestamps").BoolVar(&s.raw)
	}

	return nil
}

func (s *tSubCommand) Configure() error {
	return commonConfigure()
}

func (s *tSubCommand) Run(wg *sync.WaitGroup) (err error) {
	defer wg.Done()

	servers := func() ([]srvcache.Server, error) { return c.MiddlewareServers() }
	log := c.Logger("sub")
	conn, err := c.NewConnector(ctx, servers, c.Certname(), log)
	if err != nil {
		return fmt.Errorf("cannot connect: %s", err)
	}

	if !s.raw {
		fmt.Printf("Waiting for messages from topic %s on %s\n", s.topic, conn.ConnectedServer())
	}

	msgs := make(chan *choria.ConnectorMessage, 100)

	err = conn.QueueSubscribe(ctx, c.UniqueID(), s.topic, "", msgs)
	if err != nil {
		return fmt.Errorf("could not subscribe to %s: %s", s.topic, err)
	}

	for {
		select {
		case m := <-msgs:
			if s.raw {
				fmt.Println(string(m.Data))
				continue
			}

			if m.Subject == s.topic {
				fmt.Printf("---- %s\n%s\n\n", time.Now().Format("15:04:05"), string(m.Data))
			} else {
				fmt.Printf("---- %s on topic %s\n%s\n\n", time.Now().Format("15:04:05"), m.Subject, string(m.Data))
			}

		case <-ctx.Done():
			return nil
		}
	}
}

func init() {
	cli.commands = append(cli.commands, &tSubCommand{})
}
