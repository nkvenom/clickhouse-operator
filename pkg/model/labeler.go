// Copyright 2019 Altinity Ltd and/or its affiliates. All rights reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package model

import (
	"fmt"
	chi "github.com/altinity/clickhouse-operator/pkg/apis/clickhouse.altinity.com/v1"
	meta "k8s.io/apimachinery/pkg/apis/meta/v1"
	kublabels "k8s.io/apimachinery/pkg/labels"
)

type Labeler struct {
	version string
	chi     *chi.ClickHouseInstallation
}

func NewLabeler(version string, chi *chi.ClickHouseInstallation) *Labeler {
	return &Labeler{
		version: version,
		chi:     chi,
	}
}

func (l *Labeler) getLabelsChiScope() map[string]string {
	return map[string]string{
		LabelApp:  LabelAppValue,
		LabelChop: l.version,
		LabelChi:  getNamePartChiName(l.chi),
	}
}

func (l *Labeler) getSelectorChiScope() map[string]string {
	return map[string]string{
		LabelApp: LabelAppValue,
		// Skip chop
		LabelChi: getNamePartChiName(l.chi),
	}
}

func (l *Labeler) getLabelsClusterScope(cluster *chi.ChiCluster) map[string]string {
	return map[string]string{
		LabelApp:     LabelAppValue,
		LabelChop:    l.version,
		LabelChi:     getNamePartChiName(cluster),
		LabelCluster: getNamePartClusterName(cluster),
	}
}

func (l *Labeler) getSelectorClusterScope(cluster *chi.ChiCluster) map[string]string {
	return map[string]string{
		LabelApp: LabelAppValue,
		// Skip chop
		LabelChi:     getNamePartChiName(cluster),
		LabelCluster: getNamePartClusterName(cluster),
	}
}

func (l *Labeler) getLabelsShardScope(shard *chi.ChiShard) map[string]string {
	return map[string]string{
		LabelApp:     LabelAppValue,
		LabelChop:    l.version,
		LabelChi:     getNamePartChiName(shard),
		LabelCluster: getNamePartClusterName(shard),
		LabelShard:   getNamePartShardName(shard),
	}
}

func (l *Labeler) getSelectorShardScope(shard *chi.ChiShard) map[string]string {
	return map[string]string{
		LabelApp: LabelAppValue,
		// Skip chop
		LabelChi:     getNamePartChiName(shard),
		LabelCluster: getNamePartClusterName(shard),
		LabelShard:   getNamePartShardName(shard),
	}
}

func (l *Labeler) getLabelsHostScope(host *chi.ChiHost, applySupplementaryServiceLabels bool) map[string]string {
	labels := map[string]string{
		LabelApp:     LabelAppValue,
		LabelChop:    l.version,
		LabelChi:     getNamePartChiName(host),
		LabelCluster: getNamePartClusterName(host),
		LabelShard:   getNamePartShardName(host),
		LabelReplica: getNamePartReplicaName(host),
	}
	if applySupplementaryServiceLabels {
		labels[LabelZookeeperConfigVersion] = host.Config.ZookeeperFingerprint
		labels[LabelSettingsConfigVersion] = host.Config.SettingsFingerprint
	}
	return labels
}

func (l *Labeler) GetSelectorHostScope(host *chi.ChiHost) map[string]string {
	return map[string]string{
		LabelApp: LabelAppValue,
		// skip chop
		LabelChi:     getNamePartChiName(host),
		LabelCluster: getNamePartClusterName(host),
		LabelShard:   getNamePartShardName(host),
		LabelReplica: getNamePartReplicaName(host),
		// skip StatefulSet
		// skip Zookeeper
	}
}

func GetSelectorHostFromObjectMeta(obj *meta.ObjectMeta) (kublabels.Set, error) {
	labelApp, ok1 := obj.Labels[LabelApp]
	// skip chop
	labelChi, ok2 := obj.Labels[LabelChi]
	labelCluster, ok3 := obj.Labels[LabelCluster]
	labelShard, ok4 := obj.Labels[LabelShard]
	labelReplica, ok5 := obj.Labels[LabelReplica]
	// skip StatefulSet
	// skip Zookeeper

	if ok1 && ok2 && ok3 && ok4 && ok5 {
		set := kublabels.Set{
			LabelApp: labelApp,
			// skip chop
			LabelChi:     labelChi,
			LabelCluster: labelCluster,
			LabelShard:   labelShard,
			LabelReplica: labelReplica,
			// skip StatefulSet
			// skip Zookeeper
		}

		return set, nil
	}

	return nil, fmt.Errorf("unable to make selector from object")
}

// IsChopGeneratedObject check whether object is generated by an operator. Check is label-based
func IsChopGeneratedObject(objectMeta *meta.ObjectMeta) bool {

	// ObjectMeta must have some labels
	if len(objectMeta.Labels) == 0 {
		return false
	}

	// ObjectMeta must have LabelChop
	_, ok := objectMeta.Labels[LabelChop]

	return ok
}

func GetChiNameFromObjectMeta(meta *meta.ObjectMeta) (string, error) {
	// ObjectMeta must have LabelChi:  chi.Name label
	name, ok := meta.Labels[LabelChi]
	if ok {
		return name, nil
	} else {
		return "", fmt.Errorf("can not find %s label in meta", LabelChi)
	}
}

func GetClusterNameFromObjectMeta(meta *meta.ObjectMeta) (string, error) {
	// ObjectMeta must have LabelCluster
	name, ok := meta.Labels[LabelCluster]
	if ok {
		return name, nil
	} else {
		return "", fmt.Errorf("can not find %s label in meta", LabelChi)
	}
}
