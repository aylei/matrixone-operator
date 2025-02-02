// Copyright 2023 Matrix Origin
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package hacli

import (
	"context"
	"github.com/go-logr/logr"
	recon "github.com/matrixorigin/controller-runtime/pkg/reconciler"
	"github.com/matrixorigin/matrixone-operator/api/core/v1alpha1"
	"github.com/matrixorigin/matrixone/pkg/logservice"
	"github.com/pkg/errors"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sync"
	"time"
)

const (
	HAKeeperTimeout = 2 * time.Second
)

type HAKeeperClientManager struct {
	logger  logr.Logger
	done    chan struct{}
	kubeCli client.Client

	sync.Mutex
	logSetToClients map[types.UID]proxyClient
}

type proxyClient struct {
	client logservice.ProxyHAKeeperClient
	lsRef  types.NamespacedName
}

func NewManager(kubeCli client.Client, logger logr.Logger) *HAKeeperClientManager {
	mgr := &HAKeeperClientManager{
		logger:          logger,
		done:            make(chan struct{}),
		kubeCli:         kubeCli,
		logSetToClients: map[types.UID]proxyClient{},
	}
	go mgr.gc()
	return mgr
}

func (m *HAKeeperClientManager) GetClient(ls *v1alpha1.LogSet) (logservice.ProxyHAKeeperClient, error) {
	// FIXME: this is would be bottleneck if we concurrently reconcile a large amount of
	// matrixone clusters, we can concurrently initialize the HAKeeper clients here if necessary.
	m.Lock()
	defer m.Unlock()
	if _, ok := m.logSetToClients[ls.UID]; !ok {
		cli, err := m.newHAKeeperClient(ls)
		if err != nil {
			return nil, err
		}
		m.logSetToClients[ls.UID] = proxyClient{
			client: cli,
			lsRef:  client.ObjectKeyFromObject(ls),
		}
	}
	return m.logSetToClients[ls.UID].client, nil
}

func (m *HAKeeperClientManager) Close() {
	close(m.done)
}

func (m *HAKeeperClientManager) gc() {
	for {
		select {
		case <-time.Tick(30 * time.Second):
			m.doGC()
		case <-m.done:
			return
		}
	}
}

func (m *HAKeeperClientManager) doGC() {
	m.Lock()
	defer m.Unlock()
	for uid, v := range m.logSetToClients {
		closeFn := func() {
			err := v.client.Close()
			if err != nil {
				m.logger.Error(err, "error closing HAKeeper client", "logset", v.lsRef, "uid", uid)
			}
		}
		ls := &v1alpha1.LogSet{}
		err := m.kubeCli.Get(context.TODO(), v.lsRef, ls)
		if err != nil {
			if apierrors.IsNotFound(err) {
				delete(m.logSetToClients, uid)
				go closeFn()
			}
			m.logger.Error(err, "error gc HAKeeper client", "logset", v.lsRef, "uid", uid)
			continue
		}
		if ls.UID != uid || recon.IsReady(ls) {
			delete(m.logSetToClients, uid)
			go closeFn()
		}
	}
}

func (m *HAKeeperClientManager) newHAKeeperClient(ls *v1alpha1.LogSet) (logservice.ProxyHAKeeperClient, error) {
	ctx, cancel := context.WithTimeout(context.Background(), HAKeeperTimeout)
	defer cancel()
	cli, err := logservice.NewProxyHAKeeperClient(ctx, logservice.HAKeeperClientConfig{DiscoveryAddress: ls.Status.Discovery.String()})
	if err != nil {
		return nil, errors.Wrap(err, "build HAKeeper client")
	}
	return cli, nil
}
