<p>Packages:</p>
<ul>
<li>
<a href="#app.skpr.io">app.skpr.io</a>
</li>
<li>
<a href="#aws.skpr.io">aws.skpr.io</a>
</li>
<li>
<a href="#edge.skpr.io">edge.skpr.io</a>
</li>
<li>
<a href="#extensions.skpr.io">extensions.skpr.io</a>
</li>
<li>
<a href="#mysql.skpr.io">mysql.skpr.io</a>
</li>
</ul>
<h2 id="app.skpr.io">app.skpr.io</h2>
<p>
<p>Package v1beta1 contains API Schema definitions for the app v1beta1 API group</p>
</p>
Resource Types:
<ul><li>
<a href="#Drupal">Drupal</a>
</li></ul>
<h3 id="Drupal">Drupal
</h3>
<p>
<p>Drupal is the Schema for the drupals API</p>
</p>
<table>
<thead>
<tr>
<th>Field</th>
<th>Description</th>
</tr>
</thead>
<tbody>
<tr>
<td>
<code>apiVersion</code></br>
string</td>
<td>
<code>
app.skpr.io/v1beta1
</code>
</td>
</tr>
<tr>
<td>
<code>kind</code></br>
string
</td>
<td><code>Drupal</code></td>
</tr>
<tr>
<td>
<code>metadata</code></br>
<em>
<a href="https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.13/#objectmeta-v1-meta">
Kubernetes meta/v1.ObjectMeta
</a>
</em>
</td>
<td>
Refer to the Kubernetes API documentation for the fields of the
<code>metadata</code> field.
</td>
</tr>
<tr>
<td>
<code>spec</code></br>
<em>
<a href="#DrupalSpec">
DrupalSpec
</a>
</em>
</td>
<td>
<br/>
<br/>
<table>
<tr>
<td>
<code>nginx</code></br>
<em>
<a href="#DrupalSpecNginx">
DrupalSpecNginx
</a>
</em>
</td>
<td>
<p>Configuration for the Nginx deployment eg. image / resources / scaling.</p>
</td>
</tr>
<tr>
<td>
<code>fpm</code></br>
<em>
<a href="#DrupalSpecFPM">
DrupalSpecFPM
</a>
</em>
</td>
<td>
<p>Configuration for the FPM deployment eg. image / resources / scaling.</p>
</td>
</tr>
<tr>
<td>
<code>exec</code></br>
<em>
<a href="#DrupalSpecExec">
DrupalSpecExec
</a>
</em>
</td>
<td>
<p>Configuration for the Execution environment eg. image / resources / timeout.</p>
</td>
</tr>
<tr>
<td>
<code>volume</code></br>
<em>
<a href="#DrupalSpecVolumes">
DrupalSpecVolumes
</a>
</em>
</td>
<td>
<p>Volumes which are provisioned for the Drupal application.</p>
</td>
</tr>
<tr>
<td>
<code>mysql</code></br>
<em>
<a href="#DrupalSpecMySQL">
map[string]github.com/skpr/operator/pkg/apis/app/v1beta1.DrupalSpecMySQL
</a>
</em>
</td>
<td>
<p>Database provisioned as part of the application eg. &ldquo;default&rdquo; and &ldquo;migrate&rdquo;.</p>
</td>
</tr>
<tr>
<td>
<code>cron</code></br>
<em>
<a href="#DrupalSpecCron">
map[string]github.com/skpr/operator/pkg/apis/app/v1beta1.DrupalSpecCron
</a>
</em>
</td>
<td>
<p>Background tasks which are executed periodically eg. &ldquo;drush cron&rdquo;</p>
</td>
</tr>
<tr>
<td>
<code>configmap</code></br>
<em>
<a href="#DrupalSpecConfigMaps">
DrupalSpecConfigMaps
</a>
</em>
</td>
<td>
<p>Configuration which is exposed to the Drupal application eg. database hostname.</p>
</td>
</tr>
<tr>
<td>
<code>secret</code></br>
<em>
<a href="#DrupalSpecSecrets">
DrupalSpecSecrets
</a>
</em>
</td>
<td>
<p>Secrets which are exposed to the Drupal application eg. database credentials.</p>
</td>
</tr>
<tr>
<td>
<code>newrelic</code></br>
<em>
<a href="#DrupalSpecNewRelic">
DrupalSpecNewRelic
</a>
</em>
</td>
<td>
<p>NewRelic configuration for performance and debugging.</p>
</td>
</tr>
<tr>
<td>
<code>smtp</code></br>
<em>
<a href="#DrupalSpecSMTP">
DrupalSpecSMTP
</a>
</em>
</td>
<td>
<p>SMTP configuration for outbound email.</p>
</td>
</tr>
</table>
</td>
</tr>
<tr>
<td>
<code>status</code></br>
<em>
<a href="#DrupalStatus">
DrupalStatus
</a>
</em>
</td>
<td>
</td>
</tr>
</tbody>
</table>
<h3 id="DrupalSpec">DrupalSpec
</h3>
<p>
(<em>Appears on:</em>
<a href="#Drupal">Drupal</a>)
</p>
<p>
<p>DrupalSpec defines the desired state of Drupal</p>
</p>
<table>
<thead>
<tr>
<th>Field</th>
<th>Description</th>
</tr>
</thead>
<tbody>
<tr>
<td>
<code>nginx</code></br>
<em>
<a href="#DrupalSpecNginx">
DrupalSpecNginx
</a>
</em>
</td>
<td>
<p>Configuration for the Nginx deployment eg. image / resources / scaling.</p>
</td>
</tr>
<tr>
<td>
<code>fpm</code></br>
<em>
<a href="#DrupalSpecFPM">
DrupalSpecFPM
</a>
</em>
</td>
<td>
<p>Configuration for the FPM deployment eg. image / resources / scaling.</p>
</td>
</tr>
<tr>
<td>
<code>exec</code></br>
<em>
<a href="#DrupalSpecExec">
DrupalSpecExec
</a>
</em>
</td>
<td>
<p>Configuration for the Execution environment eg. image / resources / timeout.</p>
</td>
</tr>
<tr>
<td>
<code>volume</code></br>
<em>
<a href="#DrupalSpecVolumes">
DrupalSpecVolumes
</a>
</em>
</td>
<td>
<p>Volumes which are provisioned for the Drupal application.</p>
</td>
</tr>
<tr>
<td>
<code>mysql</code></br>
<em>
<a href="#DrupalSpecMySQL">
map[string]github.com/skpr/operator/pkg/apis/app/v1beta1.DrupalSpecMySQL
</a>
</em>
</td>
<td>
<p>Database provisioned as part of the application eg. &ldquo;default&rdquo; and &ldquo;migrate&rdquo;.</p>
</td>
</tr>
<tr>
<td>
<code>cron</code></br>
<em>
<a href="#DrupalSpecCron">
map[string]github.com/skpr/operator/pkg/apis/app/v1beta1.DrupalSpecCron
</a>
</em>
</td>
<td>
<p>Background tasks which are executed periodically eg. &ldquo;drush cron&rdquo;</p>
</td>
</tr>
<tr>
<td>
<code>configmap</code></br>
<em>
<a href="#DrupalSpecConfigMaps">
DrupalSpecConfigMaps
</a>
</em>
</td>
<td>
<p>Configuration which is exposed to the Drupal application eg. database hostname.</p>
</td>
</tr>
<tr>
<td>
<code>secret</code></br>
<em>
<a href="#DrupalSpecSecrets">
DrupalSpecSecrets
</a>
</em>
</td>
<td>
<p>Secrets which are exposed to the Drupal application eg. database credentials.</p>
</td>
</tr>
<tr>
<td>
<code>newrelic</code></br>
<em>
<a href="#DrupalSpecNewRelic">
DrupalSpecNewRelic
</a>
</em>
</td>
<td>
<p>NewRelic configuration for performance and debugging.</p>
</td>
</tr>
<tr>
<td>
<code>smtp</code></br>
<em>
<a href="#DrupalSpecSMTP">
DrupalSpecSMTP
</a>
</em>
</td>
<td>
<p>SMTP configuration for outbound email.</p>
</td>
</tr>
</tbody>
</table>
<h3 id="DrupalSpecConfigMap">DrupalSpecConfigMap
</h3>
<p>
(<em>Appears on:</em>
<a href="#DrupalSpecConfigMaps">DrupalSpecConfigMaps</a>)
</p>
<p>
<p>DrupalSpecConfigMap defines the spec for a specific configmap.</p>
</p>
<table>
<thead>
<tr>
<th>Field</th>
<th>Description</th>
</tr>
</thead>
<tbody>
<tr>
<td>
<code>path</code></br>
<em>
string
</em>
</td>
<td>
<p>Path which the configuration will be mounted.</p>
</td>
</tr>
</tbody>
</table>
<h3 id="DrupalSpecConfigMaps">DrupalSpecConfigMaps
</h3>
<p>
(<em>Appears on:</em>
<a href="#DrupalSpec">DrupalSpec</a>)
</p>
<p>
<p>DrupalSpecConfigMaps defines the spec for all configmaps.</p>
</p>
<table>
<thead>
<tr>
<th>Field</th>
<th>Description</th>
</tr>
</thead>
<tbody>
<tr>
<td>
<code>default</code></br>
<em>
<a href="#DrupalSpecConfigMap">
DrupalSpecConfigMap
</a>
</em>
</td>
<td>
<p>Configuration which is automatically set.</p>
</td>
</tr>
<tr>
<td>
<code>override</code></br>
<em>
<a href="#DrupalSpecConfigMap">
DrupalSpecConfigMap
</a>
</em>
</td>
<td>
<p>Configuration which is user provided.</p>
</td>
</tr>
</tbody>
</table>
<h3 id="DrupalSpecCron">DrupalSpecCron
</h3>
<p>
(<em>Appears on:</em>
<a href="#DrupalSpec">DrupalSpec</a>)
</p>
<p>
<p>DrupalSpecCron configures a background task for this Drupal.</p>
</p>
<table>
<thead>
<tr>
<th>Field</th>
<th>Description</th>
</tr>
</thead>
<tbody>
<tr>
<td>
<code>image</code></br>
<em>
string
</em>
</td>
<td>
<p>Image which will be executed for a background task.</p>
</td>
</tr>
<tr>
<td>
<code>readOnly</code></br>
<em>
bool
</em>
</td>
<td>
<p>Marks the filesystem as readonly to ensure code cannot be changed.</p>
</td>
</tr>
<tr>
<td>
<code>command</code></br>
<em>
string
</em>
</td>
<td>
<p>Command which will be executed for the background task.</p>
</td>
</tr>
<tr>
<td>
<code>schedule</code></br>
<em>
string
</em>
</td>
<td>
<p>Schedule whic the background task will be executed.</p>
</td>
</tr>
<tr>
<td>
<code>resources</code></br>
<em>
<a href="https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.13/#resourcerequirements-v1-core">
Kubernetes core/v1.ResourceRequirements
</a>
</em>
</td>
<td>
<p>Resource constraints which will be applied to the background task.</p>
</td>
</tr>
<tr>
<td>
<code>retries</code></br>
<em>
int32
</em>
</td>
<td>
<p>How many times to execute before marking the task as a failure.</p>
</td>
</tr>
<tr>
<td>
<code>keepSuccess</code></br>
<em>
int32
</em>
</td>
<td>
<p>How many successful builds to keep.</p>
</td>
</tr>
<tr>
<td>
<code>keepFailed</code></br>
<em>
int32
</em>
</td>
<td>
<p>How many failed builds to keep.</p>
</td>
</tr>
</tbody>
</table>
<h3 id="DrupalSpecExec">DrupalSpecExec
</h3>
<p>
(<em>Appears on:</em>
<a href="#DrupalSpec">DrupalSpec</a>)
</p>
<p>
<p>DrupalSpecExec provides configuration for an exec template.</p>
</p>
<table>
<thead>
<tr>
<th>Field</th>
<th>Description</th>
</tr>
</thead>
<tbody>
<tr>
<td>
<code>image</code></br>
<em>
string
</em>
</td>
<td>
<p>Image which will be used for developer command line access.</p>
</td>
</tr>
<tr>
<td>
<code>readOnly</code></br>
<em>
bool
</em>
</td>
<td>
<p>Marks the filesystem as readonly to ensure code cannot be changed.</p>
</td>
</tr>
<tr>
<td>
<code>resources</code></br>
<em>
<a href="https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.13/#resourcerequirements-v1-core">
Kubernetes core/v1.ResourceRequirements
</a>
</em>
</td>
<td>
<p>Resource constraints which will be applied to the background task.</p>
</td>
</tr>
<tr>
<td>
<code>timeout</code></br>
<em>
int
</em>
</td>
<td>
<p>How long for the command line environment to exist.</p>
</td>
</tr>
</tbody>
</table>
<h3 id="DrupalSpecFPM">DrupalSpecFPM
</h3>
<p>
(<em>Appears on:</em>
<a href="#DrupalSpec">DrupalSpec</a>)
</p>
<p>
<p>DrupalSpecFPM provides a specification for the FPM layer.</p>
</p>
<table>
<thead>
<tr>
<th>Field</th>
<th>Description</th>
</tr>
</thead>
<tbody>
<tr>
<td>
<code>image</code></br>
<em>
string
</em>
</td>
<td>
<p>Image which will be rolled out for the deployment.</p>
</td>
</tr>
<tr>
<td>
<code>port</code></br>
<em>
int32
</em>
</td>
<td>
<p>Port which PHP-FPM is running on.</p>
</td>
</tr>
<tr>
<td>
<code>readOnly</code></br>
<em>
bool
</em>
</td>
<td>
<p>Marks the filesystem as readonly to ensure code cannot be changed.</p>
</td>
</tr>
<tr>
<td>
<code>resources</code></br>
<em>
<a href="https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.13/#resourcerequirements-v1-core">
Kubernetes core/v1.ResourceRequirements
</a>
</em>
</td>
<td>
<p>Resource constraints which will be applied to the deployment.</p>
</td>
</tr>
<tr>
<td>
<code>autoscaling</code></br>
<em>
<a href="#DrupalSpecFPMAutoscaling">
DrupalSpecFPMAutoscaling
</a>
</em>
</td>
<td>
<p>Autoscaling rules which will be applied to the deployment.</p>
</td>
</tr>
</tbody>
</table>
<h3 id="DrupalSpecFPMAutoscaling">DrupalSpecFPMAutoscaling
</h3>
<p>
(<em>Appears on:</em>
<a href="#DrupalSpecFPM">DrupalSpecFPM</a>)
</p>
<p>
<p>DrupalSpecFPMAutoscaling contain autoscaling rules for the FPM layer.</p>
</p>
<table>
<thead>
<tr>
<th>Field</th>
<th>Description</th>
</tr>
</thead>
<tbody>
<tr>
<td>
<code>trigger</code></br>
<em>
<a href="#DrupalSpecFPMAutoscalingTrigger">
DrupalSpecFPMAutoscalingTrigger
</a>
</em>
</td>
<td>
<p>Threshold which will trigger autoscaling events.</p>
</td>
</tr>
<tr>
<td>
<code>replicas</code></br>
<em>
<a href="#DrupalSpecFPMAutoscalingReplicas">
DrupalSpecFPMAutoscalingReplicas
</a>
</em>
</td>
<td>
<p>How many copies of the deployment which should be running.</p>
</td>
</tr>
</tbody>
</table>
<h3 id="DrupalSpecFPMAutoscalingReplicas">DrupalSpecFPMAutoscalingReplicas
</h3>
<p>
(<em>Appears on:</em>
<a href="#DrupalSpecFPMAutoscaling">DrupalSpecFPMAutoscaling</a>)
</p>
<p>
<p>DrupalSpecFPMAutoscalingReplicas contain autoscaling replica rules for the FPM layer.</p>
</p>
<table>
<thead>
<tr>
<th>Field</th>
<th>Description</th>
</tr>
</thead>
<tbody>
<tr>
<td>
<code>min</code></br>
<em>
int32
</em>
</td>
<td>
<p>Minimum number of replicas for the deployment.</p>
</td>
</tr>
<tr>
<td>
<code>max</code></br>
<em>
int32
</em>
</td>
<td>
<p>Maximum number of replicas for the deployment.</p>
</td>
</tr>
</tbody>
</table>
<h3 id="DrupalSpecFPMAutoscalingTrigger">DrupalSpecFPMAutoscalingTrigger
</h3>
<p>
(<em>Appears on:</em>
<a href="#DrupalSpecFPMAutoscaling">DrupalSpecFPMAutoscaling</a>)
</p>
<p>
<p>DrupalSpecFPMAutoscalingTrigger contain autoscaling triggers for the FPM layer.</p>
</p>
<table>
<thead>
<tr>
<th>Field</th>
<th>Description</th>
</tr>
</thead>
<tbody>
<tr>
<td>
<code>cpu</code></br>
<em>
int32
</em>
</td>
<td>
<p>CPU threshold which will trigger autoscaling events.</p>
</td>
</tr>
</tbody>
</table>
<h3 id="DrupalSpecMySQL">DrupalSpecMySQL
</h3>
<p>
(<em>Appears on:</em>
<a href="#DrupalSpec">DrupalSpec</a>)
</p>
<p>
<p>DrupalSpecMySQL configuration for this Drupal.</p>
</p>
<table>
<thead>
<tr>
<th>Field</th>
<th>Description</th>
</tr>
</thead>
<tbody>
<tr>
<td>
<code>class</code></br>
<em>
string
</em>
</td>
<td>
<p>Database class which will be used when provisioning the database.</p>
</td>
</tr>
</tbody>
</table>
<h3 id="DrupalSpecNewRelic">DrupalSpecNewRelic
</h3>
<p>
(<em>Appears on:</em>
<a href="#DrupalSpec">DrupalSpec</a>)
</p>
<p>
<p>DrupalSpecNewRelic is used for configuring New Relic configuration.</p>
</p>
<table>
<thead>
<tr>
<th>Field</th>
<th>Description</th>
</tr>
</thead>
<tbody>
<tr>
<td>
<code>configmap</code></br>
<em>
<a href="#DrupalSpecNewRelicConfigMap">
DrupalSpecNewRelicConfigMap
</a>
</em>
</td>
<td>
<p>ConfigMap which contains New Relic configuration.</p>
</td>
</tr>
<tr>
<td>
<code>secret</code></br>
<em>
<a href="#DrupalSpecNewRelicSecret">
DrupalSpecNewRelicSecret
</a>
</em>
</td>
<td>
<p>Secret which contains New Relic configuration.</p>
</td>
</tr>
</tbody>
</table>
<h3 id="DrupalSpecNewRelicConfigMap">DrupalSpecNewRelicConfigMap
</h3>
<p>
(<em>Appears on:</em>
<a href="#DrupalSpecNewRelic">DrupalSpecNewRelic</a>)
</p>
<p>
<p>DrupalSpecNewRelicConfigMap for loading New Relic config from a ConfigMap.</p>
</p>
<table>
<thead>
<tr>
<th>Field</th>
<th>Description</th>
</tr>
</thead>
<tbody>
<tr>
<td>
<code>enabled</code></br>
<em>
string
</em>
</td>
<td>
<p>Key which determines if New Relic is enabled.</p>
</td>
</tr>
<tr>
<td>
<code>name</code></br>
<em>
string
</em>
</td>
<td>
<p>Name of the New Relic application.</p>
</td>
</tr>
</tbody>
</table>
<h3 id="DrupalSpecNewRelicSecret">DrupalSpecNewRelicSecret
</h3>
<p>
(<em>Appears on:</em>
<a href="#DrupalSpecNewRelic">DrupalSpecNewRelic</a>)
</p>
<p>
<p>DrupalSpecNewRelicSecret for loading New Relic config from a Secret.</p>
</p>
<table>
<thead>
<tr>
<th>Field</th>
<th>Description</th>
</tr>
</thead>
<tbody>
<tr>
<td>
<code>license</code></br>
<em>
string
</em>
</td>
<td>
<p>License (API Key) for interacting with New Relic.</p>
</td>
</tr>
</tbody>
</table>
<h3 id="DrupalSpecNginx">DrupalSpecNginx
</h3>
<p>
(<em>Appears on:</em>
<a href="#DrupalSpec">DrupalSpec</a>)
</p>
<p>
<p>DrupalSpecNginx provides a specification for the Nginx layer.</p>
</p>
<table>
<thead>
<tr>
<th>Field</th>
<th>Description</th>
</tr>
</thead>
<tbody>
<tr>
<td>
<code>image</code></br>
<em>
string
</em>
</td>
<td>
<p>Image which will be rolled out for the deployment.</p>
</td>
</tr>
<tr>
<td>
<code>port</code></br>
<em>
int32
</em>
</td>
<td>
<p>Port which Nginx is running on.</p>
</td>
</tr>
<tr>
<td>
<code>readOnly</code></br>
<em>
bool
</em>
</td>
<td>
<p>Marks the filesystem as readonly to ensure code cannot be changed.</p>
</td>
</tr>
<tr>
<td>
<code>resources</code></br>
<em>
<a href="https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.13/#resourcerequirements-v1-core">
Kubernetes core/v1.ResourceRequirements
</a>
</em>
</td>
<td>
<p>Resource constraints which will be applied to the deployment.</p>
</td>
</tr>
<tr>
<td>
<code>autoscaling</code></br>
<em>
<a href="#DrupalSpecNginxAutoscaling">
DrupalSpecNginxAutoscaling
</a>
</em>
</td>
<td>
<p>Autoscaling rules which will be applied to the deployment.</p>
</td>
</tr>
<tr>
<td>
<code>hostAlias</code></br>
<em>
<a href="#DrupalSpecNginxHostAlias">
DrupalSpecNginxHostAlias
</a>
</em>
</td>
<td>
<p>HostAlias which connects Nginx to the FPM service.</p>
</td>
</tr>
</tbody>
</table>
<h3 id="DrupalSpecNginxAutoscaling">DrupalSpecNginxAutoscaling
</h3>
<p>
(<em>Appears on:</em>
<a href="#DrupalSpecNginx">DrupalSpecNginx</a>)
</p>
<p>
<p>DrupalSpecNginxAutoscaling contain autoscaling rules for the Nginx layer.</p>
</p>
<table>
<thead>
<tr>
<th>Field</th>
<th>Description</th>
</tr>
</thead>
<tbody>
<tr>
<td>
<code>trigger</code></br>
<em>
<a href="#DrupalSpecNginxAutoscalingTrigger">
DrupalSpecNginxAutoscalingTrigger
</a>
</em>
</td>
<td>
<p>Threshold which will trigger autoscaling events.</p>
</td>
</tr>
<tr>
<td>
<code>replicas</code></br>
<em>
<a href="#DrupalSpecNginxAutoscalingReplicas">
DrupalSpecNginxAutoscalingReplicas
</a>
</em>
</td>
<td>
<p>How many copies of the deployment which should be running.</p>
</td>
</tr>
</tbody>
</table>
<h3 id="DrupalSpecNginxAutoscalingReplicas">DrupalSpecNginxAutoscalingReplicas
</h3>
<p>
(<em>Appears on:</em>
<a href="#DrupalSpecNginxAutoscaling">DrupalSpecNginxAutoscaling</a>)
</p>
<p>
<p>DrupalSpecNginxAutoscalingReplicas contain autoscaling replica rules for the Nginx layer.</p>
</p>
<table>
<thead>
<tr>
<th>Field</th>
<th>Description</th>
</tr>
</thead>
<tbody>
<tr>
<td>
<code>min</code></br>
<em>
int32
</em>
</td>
<td>
<p>Minimum number of replicas for the deployment.</p>
</td>
</tr>
<tr>
<td>
<code>max</code></br>
<em>
int32
</em>
</td>
<td>
<p>Maximum number of replicas for the deployment.</p>
</td>
</tr>
</tbody>
</table>
<h3 id="DrupalSpecNginxAutoscalingTrigger">DrupalSpecNginxAutoscalingTrigger
</h3>
<p>
(<em>Appears on:</em>
<a href="#DrupalSpecNginxAutoscaling">DrupalSpecNginxAutoscaling</a>)
</p>
<p>
<p>DrupalSpecNginxAutoscalingTrigger contain autoscaling triggers for the Nginx layer.</p>
</p>
<table>
<thead>
<tr>
<th>Field</th>
<th>Description</th>
</tr>
</thead>
<tbody>
<tr>
<td>
<code>cpu</code></br>
<em>
int32
</em>
</td>
<td>
<p>CPU threshold which will trigger autoscaling events.</p>
</td>
</tr>
</tbody>
</table>
<h3 id="DrupalSpecNginxHostAlias">DrupalSpecNginxHostAlias
</h3>
<p>
(<em>Appears on:</em>
<a href="#DrupalSpecNginx">DrupalSpecNginx</a>)
</p>
<p>
<p>DrupalSpecNginxHostAlias declares static DNS entries which wire up the Nginx layer to the FPM lay</p>
</p>
<table>
<thead>
<tr>
<th>Field</th>
<th>Description</th>
</tr>
</thead>
<tbody>
<tr>
<td>
<code>fpm</code></br>
<em>
string
</em>
</td>
<td>
<p>HostAlias for the PHP FPM service.</p>
</td>
</tr>
</tbody>
</table>
<h3 id="DrupalSpecSMTP">DrupalSpecSMTP
</h3>
<p>
(<em>Appears on:</em>
<a href="#DrupalSpec">DrupalSpec</a>)
</p>
<p>
<p>DrupalSpecSMTP configuration for outbound email.</p>
</p>
<table>
<thead>
<tr>
<th>Field</th>
<th>Description</th>
</tr>
</thead>
<tbody>
<tr>
<td>
<code>from</code></br>
<em>
<a href="#SMTPSpecFrom">
SMTPSpecFrom
</a>
</em>
</td>
<td>
<p>Configuration for validationing FROM addresses.</p>
</td>
</tr>
</tbody>
</table>
<h3 id="DrupalSpecSecret">DrupalSpecSecret
</h3>
<p>
(<em>Appears on:</em>
<a href="#DrupalSpecSecrets">DrupalSpecSecrets</a>)
</p>
<p>
<p>DrupalSpecSecret defines the spec for a specific secret.</p>
</p>
<table>
<thead>
<tr>
<th>Field</th>
<th>Description</th>
</tr>
</thead>
<tbody>
<tr>
<td>
<code>path</code></br>
<em>
string
</em>
</td>
<td>
<p>Path which the secrets will be mounted.</p>
</td>
</tr>
</tbody>
</table>
<h3 id="DrupalSpecSecrets">DrupalSpecSecrets
</h3>
<p>
(<em>Appears on:</em>
<a href="#DrupalSpec">DrupalSpec</a>)
</p>
<p>
<p>DrupalSpecSecrets defines the spec for all secrets.</p>
</p>
<table>
<thead>
<tr>
<th>Field</th>
<th>Description</th>
</tr>
</thead>
<tbody>
<tr>
<td>
<code>default</code></br>
<em>
<a href="#DrupalSpecSecret">
DrupalSpecSecret
</a>
</em>
</td>
<td>
<p>Secrets which are automatically set.</p>
</td>
</tr>
<tr>
<td>
<code>override</code></br>
<em>
<a href="#DrupalSpecSecret">
DrupalSpecSecret
</a>
</em>
</td>
<td>
<p>Secrets which are user provided.</p>
</td>
</tr>
</tbody>
</table>
<h3 id="DrupalSpecVolume">DrupalSpecVolume
</h3>
<p>
(<em>Appears on:</em>
<a href="#DrupalSpecVolumes">DrupalSpecVolumes</a>)
</p>
<p>
<p>DrupalSpecVolume which will be mounted for this Drupal.</p>
</p>
<table>
<thead>
<tr>
<th>Field</th>
<th>Description</th>
</tr>
</thead>
<tbody>
<tr>
<td>
<code>path</code></br>
<em>
string
</em>
</td>
<td>
<p>Path which the volume will be mounted.</p>
</td>
</tr>
<tr>
<td>
<code>class</code></br>
<em>
string
</em>
</td>
<td>
<p>StorageClass which will be used to provision storage.</p>
</td>
</tr>
<tr>
<td>
<code>amount</code></br>
<em>
string
</em>
</td>
<td>
<p>Amount of storage which will be provisioned.</p>
</td>
</tr>
<tr>
<td>
<code>permissions</code></br>
<em>
<a href="#DrupalSpecVolumePermissions">
DrupalSpecVolumePermissions
</a>
</em>
</td>
<td>
<p>Permmissions which will be enforced for this volume.</p>
</td>
</tr>
</tbody>
</table>
<h3 id="DrupalSpecVolumePermissions">DrupalSpecVolumePermissions
</h3>
<p>
(<em>Appears on:</em>
<a href="#DrupalSpecVolume">DrupalSpecVolume</a>)
</p>
<p>
<p>DrupalSpecVolumePermissions describes permisions which will be enforced for a volume.</p>
</p>
<table>
<thead>
<tr>
<th>Field</th>
<th>Description</th>
</tr>
</thead>
<tbody>
<tr>
<td>
<code>user</code></br>
<em>
string
</em>
</td>
<td>
<p>User name which will be enforced for the volume. Can also be a UID.</p>
</td>
</tr>
<tr>
<td>
<code>group</code></br>
<em>
string
</em>
</td>
<td>
<p>Group name which will be enforced for the volume. Can also be a GID.</p>
</td>
</tr>
<tr>
<td>
<code>directory</code></br>
<em>
int
</em>
</td>
<td>
<p>Directory chmod which will be applied to the volume.</p>
</td>
</tr>
<tr>
<td>
<code>file</code></br>
<em>
int
</em>
</td>
<td>
<p>File chmod which will be applied to the volume.</p>
</td>
</tr>
<tr>
<td>
<code>cronJob</code></br>
<em>
<a href="#DrupalSpecVolumePermissionsCronJob">
DrupalSpecVolumePermissionsCronJob
</a>
</em>
</td>
<td>
<p>CronJob which will be run to enforce permissions.</p>
</td>
</tr>
</tbody>
</table>
<h3 id="DrupalSpecVolumePermissionsCronJob">DrupalSpecVolumePermissionsCronJob
</h3>
<p>
(<em>Appears on:</em>
<a href="#DrupalSpecVolumePermissions">DrupalSpecVolumePermissions</a>)
</p>
<p>
<p>DrupalSpecVolumePermissionsCronJob describes how the permissions will be enforced.</p>
</p>
<table>
<thead>
<tr>
<th>Field</th>
<th>Description</th>
</tr>
</thead>
<tbody>
<tr>
<td>
<code>image</code></br>
<em>
string
</em>
</td>
<td>
<p>Privileged image which will be used when running chown and chmod commands. Needs be an image which runs as root to enforce permissions.</p>
</td>
</tr>
<tr>
<td>
<code>schedule</code></br>
<em>
string
</em>
</td>
<td>
<p>Schedule which permissions will be checked.</p>
</td>
</tr>
<tr>
<td>
<code>deadline</code></br>
<em>
int64
</em>
</td>
<td>
<p>How long before a Job should be marked as failed if it does not get scheduled in time.</p>
</td>
</tr>
<tr>
<td>
<code>retries</code></br>
<em>
int32
</em>
</td>
<td>
<p>How many times to get a &ldquo;Successful&rdquo; execution before failing.</p>
</td>
</tr>
<tr>
<td>
<code>keepSuccess</code></br>
<em>
int32
</em>
</td>
<td>
<p>How many past successes to keep.</p>
</td>
</tr>
<tr>
<td>
<code>keepFailed</code></br>
<em>
int32
</em>
</td>
<td>
<p>How many past failures to keep.</p>
</td>
</tr>
<tr>
<td>
<code>readOnly</code></br>
<em>
bool
</em>
</td>
<td>
<p>Make the root filesystem of this image ready only.</p>
</td>
</tr>
</tbody>
</table>
<h3 id="DrupalSpecVolumes">DrupalSpecVolumes
</h3>
<p>
(<em>Appears on:</em>
<a href="#DrupalSpec">DrupalSpec</a>)
</p>
<p>
<p>DrupalSpecVolumes which will be mounted for this Drupal.</p>
</p>
<table>
<thead>
<tr>
<th>Field</th>
<th>Description</th>
</tr>
</thead>
<tbody>
<tr>
<td>
<code>public</code></br>
<em>
<a href="#DrupalSpecVolume">
DrupalSpecVolume
</a>
</em>
</td>
<td>
<p>Public filesystem which is mounted into the Nginx and PHP-FPM deployments.</p>
</td>
</tr>
<tr>
<td>
<code>private</code></br>
<em>
<a href="#DrupalSpecVolume">
DrupalSpecVolume
</a>
</em>
</td>
<td>
<p>Private filesystem which is only mounted into the PHP-FPM deployment.</p>
</td>
</tr>
<tr>
<td>
<code>temporary</code></br>
<em>
<a href="#DrupalSpecVolume">
DrupalSpecVolume
</a>
</em>
</td>
<td>
<p>Temporary filesystem which is only mounted into the PHP-FPM deployment.</p>
</td>
</tr>
</tbody>
</table>
<h3 id="DrupalStatus">DrupalStatus
</h3>
<p>
(<em>Appears on:</em>
<a href="#Drupal">Drupal</a>)
</p>
<p>
<p>DrupalStatus defines the observed state of Drupal</p>
</p>
<table>
<thead>
<tr>
<th>Field</th>
<th>Description</th>
</tr>
</thead>
<tbody>
<tr>
<td>
<code>observedGeneration</code></br>
<em>
int64
</em>
</td>
<td>
<p>Used for determining if an APIs information is up to date.</p>
</td>
</tr>
<tr>
<td>
<code>labels</code></br>
<em>
<a href="#DrupalStatusLabels">
DrupalStatusLabels
</a>
</em>
</td>
<td>
<p>Labels for querying Drupal services.</p>
</td>
</tr>
<tr>
<td>
<code>nginx</code></br>
<em>
<a href="#DrupalStatusNginx">
DrupalStatusNginx
</a>
</em>
</td>
<td>
<p>Status for the Nginx deployment.</p>
</td>
</tr>
<tr>
<td>
<code>fpm</code></br>
<em>
<a href="#DrupalStatusFPM">
DrupalStatusFPM
</a>
</em>
</td>
<td>
<p>Status for the PHP-FPM deployment.</p>
</td>
</tr>
<tr>
<td>
<code>volume</code></br>
<em>
<a href="#DrupalStatusVolumes">
DrupalStatusVolumes
</a>
</em>
</td>
<td>
<p>Volume status information.</p>
</td>
</tr>
<tr>
<td>
<code>mysql</code></br>
<em>
<a href="#DrupalStatusMySQL">
map[string]github.com/skpr/operator/pkg/apis/app/v1beta1.DrupalStatusMySQL
</a>
</em>
</td>
<td>
<p>MySQL status information.</p>
</td>
</tr>
<tr>
<td>
<code>cron</code></br>
<em>
<a href="#DrupalStatusCron">
map[string]github.com/skpr/operator/pkg/apis/app/v1beta1.DrupalStatusCron
</a>
</em>
</td>
<td>
<p>Background task information eg. Last execution.</p>
</td>
</tr>
<tr>
<td>
<code>configmap</code></br>
<em>
<a href="#DrupalStatusConfigMaps">
DrupalStatusConfigMaps
</a>
</em>
</td>
<td>
<p>Configuration status.</p>
</td>
</tr>
<tr>
<td>
<code>secret</code></br>
<em>
<a href="#DrupalStatusSecrets">
DrupalStatusSecrets
</a>
</em>
</td>
<td>
<p>Secrets status.</p>
</td>
</tr>
<tr>
<td>
<code>exec</code></br>
<em>
<a href="#DrupalStatusExec">
DrupalStatusExec
</a>
</em>
</td>
<td>
<p>Execution environment status.</p>
</td>
</tr>
<tr>
<td>
<code>smtp</code></br>
<em>
<a href="#DrupalStatusSMTP">
DrupalStatusSMTP
</a>
</em>
</td>
<td>
<p>SMTP service status.</p>
</td>
</tr>
</tbody>
</table>
<h3 id="DrupalStatusConfigMap">DrupalStatusConfigMap
</h3>
<p>
(<em>Appears on:</em>
<a href="#DrupalStatusConfigMaps">DrupalStatusConfigMaps</a>)
</p>
<p>
<p>DrupalStatusConfigMap identifies specific config map related status.</p>
</p>
<table>
<thead>
<tr>
<th>Field</th>
<th>Description</th>
</tr>
</thead>
<tbody>
<tr>
<td>
<code>name</code></br>
<em>
string
</em>
</td>
<td>
<p>Name of the configmap.</p>
</td>
</tr>
<tr>
<td>
<code>count</code></br>
<em>
int
</em>
</td>
<td>
<p>How many configuration are assigned to the application.</p>
</td>
</tr>
</tbody>
</table>
<h3 id="DrupalStatusConfigMaps">DrupalStatusConfigMaps
</h3>
<p>
(<em>Appears on:</em>
<a href="#DrupalStatus">DrupalStatus</a>)
</p>
<p>
<p>DrupalStatusConfigMaps identifies all config map related status.</p>
</p>
<table>
<thead>
<tr>
<th>Field</th>
<th>Description</th>
</tr>
</thead>
<tbody>
<tr>
<td>
<code>default</code></br>
<em>
<a href="#DrupalStatusConfigMap">
DrupalStatusConfigMap
</a>
</em>
</td>
<td>
<p>Status of the generated configuration.</p>
</td>
</tr>
<tr>
<td>
<code>override</code></br>
<em>
<a href="#DrupalStatusConfigMap">
DrupalStatusConfigMap
</a>
</em>
</td>
<td>
<p>Status of the user provided configuration.</p>
</td>
</tr>
</tbody>
</table>
<h3 id="DrupalStatusCron">DrupalStatusCron
</h3>
<p>
(<em>Appears on:</em>
<a href="#DrupalStatus">DrupalStatus</a>)
</p>
<p>
<p>DrupalStatusCron identifies all cron related status.</p>
</p>
<table>
<thead>
<tr>
<th>Field</th>
<th>Description</th>
</tr>
</thead>
<tbody>
<tr>
<td>
<code>lastScheduleTime</code></br>
<em>
<a href="https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.13/#time-v1-meta">
Kubernetes meta/v1.Time
</a>
</em>
</td>
<td>
<p>Last time a background task was executed.</p>
</td>
</tr>
</tbody>
</table>
<h3 id="DrupalStatusExec">DrupalStatusExec
</h3>
<p>
(<em>Appears on:</em>
<a href="#DrupalStatus">DrupalStatus</a>)
</p>
<p>
<p>DrupalStatusExec identifies specific template which can be loaded by external sources.</p>
</p>
<table>
<thead>
<tr>
<th>Field</th>
<th>Description</th>
</tr>
</thead>
<tbody>
<tr>
<td>
<code>name</code></br>
<em>
string
</em>
</td>
<td>
<p>Name of the command line environment template.</p>
</td>
</tr>
</tbody>
</table>
<h3 id="DrupalStatusFPM">DrupalStatusFPM
</h3>
<p>
(<em>Appears on:</em>
<a href="#DrupalStatus">DrupalStatus</a>)
</p>
<p>
<p>DrupalStatusFPM identifies all deployment related status.</p>
</p>
<table>
<thead>
<tr>
<th>Field</th>
<th>Description</th>
</tr>
</thead>
<tbody>
<tr>
<td>
<code>phase</code></br>
<em>
github.com/skpr/operator/pkg/utils/k8s/deployment.Phase
</em>
</td>
<td>
<p>Current phase of the PHP-FPM deployment eg. InProgress.</p>
</td>
</tr>
<tr>
<td>
<code>service</code></br>
<em>
string
</em>
</td>
<td>
<p>Service for routing traffic.</p>
</td>
</tr>
<tr>
<td>
<code>image</code></br>
<em>
string
</em>
</td>
<td>
<p>Current image which has been rolled out.</p>
</td>
</tr>
<tr>
<td>
<code>replicas</code></br>
<em>
int32
</em>
</td>
<td>
<p>Current number of replicas.</p>
</td>
</tr>
<tr>
<td>
<code>metrics</code></br>
<em>
<a href="#DrupalStatusFPMMetrics">
DrupalStatusFPMMetrics
</a>
</em>
</td>
<td>
<p>Application metrics.</p>
</td>
</tr>
</tbody>
</table>
<h3 id="DrupalStatusFPMMetrics">DrupalStatusFPMMetrics
</h3>
<p>
(<em>Appears on:</em>
<a href="#DrupalStatusFPM">DrupalStatusFPM</a>)
</p>
<p>
<p>DrupalStatusFPMMetrics identifies all fpm metric related status.</p>
</p>
<table>
<thead>
<tr>
<th>Field</th>
<th>Description</th>
</tr>
</thead>
<tbody>
<tr>
<td>
<code>cpu</code></br>
<em>
int32
</em>
</td>
<td>
<p>Current CPU metric for the PHP-FPM deployment.</p>
</td>
</tr>
</tbody>
</table>
<h3 id="DrupalStatusLabels">DrupalStatusLabels
</h3>
<p>
(<em>Appears on:</em>
<a href="#DrupalStatus">DrupalStatus</a>)
</p>
<p>
<p>DrupalStatusLabels which are used for querying application components.</p>
</p>
<table>
<thead>
<tr>
<th>Field</th>
<th>Description</th>
</tr>
</thead>
<tbody>
<tr>
<td>
<code>all</code></br>
<em>
map[string]string
</em>
</td>
<td>
<p>Used to query all Drupal application pods.</p>
</td>
</tr>
<tr>
<td>
<code>nginx</code></br>
<em>
map[string]string
</em>
</td>
<td>
<p>Used to query all Nginx pods.</p>
</td>
</tr>
<tr>
<td>
<code>fpm</code></br>
<em>
map[string]string
</em>
</td>
<td>
<p>Used to query all PHP-FPM pods.</p>
</td>
</tr>
<tr>
<td>
<code>cron</code></br>
<em>
map[string]string
</em>
</td>
<td>
<p>Used to query all background task pods.</p>
</td>
</tr>
</tbody>
</table>
<h3 id="DrupalStatusMySQL">DrupalStatusMySQL
</h3>
<p>
(<em>Appears on:</em>
<a href="#DrupalStatus">DrupalStatus</a>)
</p>
<p>
<p>DrupalStatusMySQL identifies all mysql related status.</p>
</p>
<table>
<thead>
<tr>
<th>Field</th>
<th>Description</th>
</tr>
</thead>
<tbody>
<tr>
<td>
<code>configmap</code></br>
<em>
<a href="#DrupalStatusMySQLConfigMap">
DrupalStatusMySQLConfigMap
</a>
</em>
</td>
<td>
<p>Status of the application configuration.</p>
</td>
</tr>
<tr>
<td>
<code>secret</code></br>
<em>
<a href="#DrupalStatusMySQLSecret">
DrupalStatusMySQLSecret
</a>
</em>
</td>
<td>
<p>Status of the application secrets.</p>
</td>
</tr>
</tbody>
</table>
<h3 id="DrupalStatusMySQLConfigMap">DrupalStatusMySQLConfigMap
</h3>
<p>
(<em>Appears on:</em>
<a href="#DrupalStatusMySQL">DrupalStatusMySQL</a>)
</p>
<p>
<p>DrupalStatusMySQLConfigMap identifies all mysql configmap related status.</p>
</p>
<table>
<thead>
<tr>
<th>Field</th>
<th>Description</th>
</tr>
</thead>
<tbody>
<tr>
<td>
<code>name</code></br>
<em>
string
</em>
</td>
<td>
<p>Name of the configmap.</p>
</td>
</tr>
<tr>
<td>
<code>keys</code></br>
<em>
<a href="#DrupalStatusMySQLConfigMapKeys">
DrupalStatusMySQLConfigMapKeys
</a>
</em>
</td>
<td>
<p>Keys which can be used for discovery.</p>
</td>
</tr>
</tbody>
</table>
<h3 id="DrupalStatusMySQLConfigMapKeys">DrupalStatusMySQLConfigMapKeys
</h3>
<p>
(<em>Appears on:</em>
<a href="#DrupalStatusMySQLConfigMap">DrupalStatusMySQLConfigMap</a>)
</p>
<p>
<p>DrupalStatusMySQLConfigMapKeys identifies all mysql configmap keys.</p>
</p>
<table>
<thead>
<tr>
<th>Field</th>
<th>Description</th>
</tr>
</thead>
<tbody>
<tr>
<td>
<code>database</code></br>
<em>
string
</em>
</td>
<td>
<p>Key which was applied to the application for database connectivity.</p>
</td>
</tr>
<tr>
<td>
<code>hostname</code></br>
<em>
string
</em>
</td>
<td>
<p>Key which was applied to the application for database connectivity.</p>
</td>
</tr>
<tr>
<td>
<code>port</code></br>
<em>
string
</em>
</td>
<td>
<p>Key which was applied to the application for database connectivity.</p>
</td>
</tr>
</tbody>
</table>
<h3 id="DrupalStatusMySQLSecret">DrupalStatusMySQLSecret
</h3>
<p>
(<em>Appears on:</em>
<a href="#DrupalStatusMySQL">DrupalStatusMySQL</a>)
</p>
<p>
<p>DrupalStatusMySQLSecret identifies all mysql secret related status.</p>
</p>
<table>
<thead>
<tr>
<th>Field</th>
<th>Description</th>
</tr>
</thead>
<tbody>
<tr>
<td>
<code>name</code></br>
<em>
string
</em>
</td>
<td>
<p>Name of the secret.</p>
</td>
</tr>
<tr>
<td>
<code>keys</code></br>
<em>
<a href="#DrupalStatusMySQLSecretKeys">
DrupalStatusMySQLSecretKeys
</a>
</em>
</td>
<td>
<p>Keys which can be used for discovery.</p>
</td>
</tr>
</tbody>
</table>
<h3 id="DrupalStatusMySQLSecretKeys">DrupalStatusMySQLSecretKeys
</h3>
<p>
(<em>Appears on:</em>
<a href="#DrupalStatusMySQLSecret">DrupalStatusMySQLSecret</a>)
</p>
<p>
<p>DrupalStatusMySQLSecretKeys identifies all mysql secret keys.</p>
</p>
<table>
<thead>
<tr>
<th>Field</th>
<th>Description</th>
</tr>
</thead>
<tbody>
<tr>
<td>
<code>username</code></br>
<em>
string
</em>
</td>
<td>
<p>Key which was applied to the application for database connectivity.</p>
</td>
</tr>
<tr>
<td>
<code>password</code></br>
<em>
string
</em>
</td>
<td>
<p>Key which was applied to the application for database connectivity.</p>
</td>
</tr>
</tbody>
</table>
<h3 id="DrupalStatusNginx">DrupalStatusNginx
</h3>
<p>
(<em>Appears on:</em>
<a href="#DrupalStatus">DrupalStatus</a>)
</p>
<p>
<p>DrupalStatusNginx identifies all deployment related status.</p>
</p>
<table>
<thead>
<tr>
<th>Field</th>
<th>Description</th>
</tr>
</thead>
<tbody>
<tr>
<td>
<code>phase</code></br>
<em>
github.com/skpr/operator/pkg/utils/k8s/deployment.Phase
</em>
</td>
<td>
<p>Current phase of the Nginx deployment eg. InProgress.</p>
</td>
</tr>
<tr>
<td>
<code>service</code></br>
<em>
string
</em>
</td>
<td>
<p>Service for routing traffic.</p>
</td>
</tr>
<tr>
<td>
<code>image</code></br>
<em>
string
</em>
</td>
<td>
<p>Current image which has been rolled out.</p>
</td>
</tr>
<tr>
<td>
<code>replicas</code></br>
<em>
int32
</em>
</td>
<td>
<p>Current number of replicas.</p>
</td>
</tr>
<tr>
<td>
<code>metrics</code></br>
<em>
<a href="#DrupalStatusNginxMetrics">
DrupalStatusNginxMetrics
</a>
</em>
</td>
<td>
<p>Application metrics.</p>
</td>
</tr>
</tbody>
</table>
<h3 id="DrupalStatusNginxMetrics">DrupalStatusNginxMetrics
</h3>
<p>
(<em>Appears on:</em>
<a href="#DrupalStatusNginx">DrupalStatusNginx</a>)
</p>
<p>
<p>DrupalStatusNginxMetrics identifies all nginx metric related status.</p>
</p>
<table>
<thead>
<tr>
<th>Field</th>
<th>Description</th>
</tr>
</thead>
<tbody>
<tr>
<td>
<code>cpu</code></br>
<em>
int32
</em>
</td>
<td>
<p>Current CPU metric for the Nginx deployment.</p>
</td>
</tr>
</tbody>
</table>
<h3 id="DrupalStatusSMTP">DrupalStatusSMTP
</h3>
<p>
(<em>Appears on:</em>
<a href="#DrupalStatus">DrupalStatus</a>)
</p>
<p>
<p>DrupalStatusSMTP provides the status for the SMTP service.</p>
</p>
<table>
<thead>
<tr>
<th>Field</th>
<th>Description</th>
</tr>
</thead>
<tbody>
<tr>
<td>
<code>verification</code></br>
<em>
<a href="#SMTPStatusVerification">
SMTPStatusVerification
</a>
</em>
</td>
<td>
</td>
</tr>
</tbody>
</table>
<h3 id="DrupalStatusSecret">DrupalStatusSecret
</h3>
<p>
(<em>Appears on:</em>
<a href="#DrupalStatusSecrets">DrupalStatusSecrets</a>)
</p>
<p>
<p>DrupalStatusSecret identifies specific secret related status.</p>
</p>
<table>
<thead>
<tr>
<th>Field</th>
<th>Description</th>
</tr>
</thead>
<tbody>
<tr>
<td>
<code>name</code></br>
<em>
string
</em>
</td>
<td>
<p>Name of the secret.</p>
</td>
</tr>
<tr>
<td>
<code>count</code></br>
<em>
int
</em>
</td>
<td>
<p>How many secrets are assigned to the application.</p>
</td>
</tr>
</tbody>
</table>
<h3 id="DrupalStatusSecrets">DrupalStatusSecrets
</h3>
<p>
(<em>Appears on:</em>
<a href="#DrupalStatus">DrupalStatus</a>)
</p>
<p>
<p>DrupalStatusSecrets identifies all secret related status.</p>
</p>
<table>
<thead>
<tr>
<th>Field</th>
<th>Description</th>
</tr>
</thead>
<tbody>
<tr>
<td>
<code>default</code></br>
<em>
<a href="#DrupalStatusSecret">
DrupalStatusSecret
</a>
</em>
</td>
<td>
<p>Status of the generated secrets.</p>
</td>
</tr>
<tr>
<td>
<code>override</code></br>
<em>
<a href="#DrupalStatusSecret">
DrupalStatusSecret
</a>
</em>
</td>
<td>
<p>Status of the user provided secrets.</p>
</td>
</tr>
</tbody>
</table>
<h3 id="DrupalStatusVolume">DrupalStatusVolume
</h3>
<p>
(<em>Appears on:</em>
<a href="#DrupalStatusVolumes">DrupalStatusVolumes</a>)
</p>
<p>
<p>DrupalStatusVolume identifies specific volume related status.</p>
</p>
<table>
<thead>
<tr>
<th>Field</th>
<th>Description</th>
</tr>
</thead>
<tbody>
<tr>
<td>
<code>name</code></br>
<em>
string
</em>
</td>
<td>
<p>Name of the volume.</p>
</td>
</tr>
<tr>
<td>
<code>phase</code></br>
<em>
<a href="https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.13/#persistentvolumeclaimphase-v1-core">
Kubernetes core/v1.PersistentVolumeClaimPhase
</a>
</em>
</td>
<td>
<p>Current state of the volume.</p>
</td>
</tr>
</tbody>
</table>
<h3 id="DrupalStatusVolumes">DrupalStatusVolumes
</h3>
<p>
(<em>Appears on:</em>
<a href="#DrupalStatus">DrupalStatus</a>)
</p>
<p>
<p>DrupalStatusVolumes identifies all volume related status.</p>
</p>
<table>
<thead>
<tr>
<th>Field</th>
<th>Description</th>
</tr>
</thead>
<tbody>
<tr>
<td>
<code>public</code></br>
<em>
<a href="#DrupalStatusVolume">
DrupalStatusVolume
</a>
</em>
</td>
<td>
<p>Current state of the public filesystem volume.</p>
</td>
</tr>
<tr>
<td>
<code>private</code></br>
<em>
<a href="#DrupalStatusVolume">
DrupalStatusVolume
</a>
</em>
</td>
<td>
<p>Current state of the private filesystem volume.</p>
</td>
</tr>
<tr>
<td>
<code>temporary</code></br>
<em>
<a href="#DrupalStatusVolume">
DrupalStatusVolume
</a>
</em>
</td>
<td>
<p>Current state of the temporary filesystem volume.</p>
</td>
</tr>
</tbody>
</table>
<hr/>
<h2 id="aws.skpr.io">aws.skpr.io</h2>
<p>
<p>Package v1beta1 contains API Schema definitions for the aws v1beta1 API group</p>
</p>
Resource Types:
<ul><li>
<a href="#Certificate">Certificate</a>
</li><li>
<a href="#CertificateRequest">CertificateRequest</a>
</li><li>
<a href="#CloudFront">CloudFront</a>
</li><li>
<a href="#CloudFrontInvalidation">CloudFrontInvalidation</a>
</li></ul>
<h3 id="Certificate">Certificate
</h3>
<p>
<p>Certificate is the Schema for the certificates API</p>
</p>
<table>
<thead>
<tr>
<th>Field</th>
<th>Description</th>
</tr>
</thead>
<tbody>
<tr>
<td>
<code>apiVersion</code></br>
string</td>
<td>
<code>
aws.skpr.io/v1beta1
</code>
</td>
</tr>
<tr>
<td>
<code>kind</code></br>
string
</td>
<td><code>Certificate</code></td>
</tr>
<tr>
<td>
<code>metadata</code></br>
<em>
<a href="https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.13/#objectmeta-v1-meta">
Kubernetes meta/v1.ObjectMeta
</a>
</em>
</td>
<td>
Refer to the Kubernetes API documentation for the fields of the
<code>metadata</code> field.
</td>
</tr>
<tr>
<td>
<code>spec</code></br>
<em>
<a href="#CertificateSpec">
CertificateSpec
</a>
</em>
</td>
<td>
<br/>
<br/>
<table>
<tr>
<td>
<code>request</code></br>
<em>
<a href="#CertificateRequestSpec">
CertificateRequestSpec
</a>
</em>
</td>
<td>
<p>Information which will be used to provision a certificate.</p>
</td>
</tr>
</table>
</td>
</tr>
<tr>
<td>
<code>status</code></br>
<em>
<a href="#CertificateStatus">
CertificateStatus
</a>
</em>
</td>
<td>
</td>
</tr>
</tbody>
</table>
<h3 id="CertificateRequest">CertificateRequest
</h3>
<p>
<p>CertificateRequest is the Schema for the certificaterequests API</p>
</p>
<table>
<thead>
<tr>
<th>Field</th>
<th>Description</th>
</tr>
</thead>
<tbody>
<tr>
<td>
<code>apiVersion</code></br>
string</td>
<td>
<code>
aws.skpr.io/v1beta1
</code>
</td>
</tr>
<tr>
<td>
<code>kind</code></br>
string
</td>
<td><code>CertificateRequest</code></td>
</tr>
<tr>
<td>
<code>metadata</code></br>
<em>
<a href="https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.13/#objectmeta-v1-meta">
Kubernetes meta/v1.ObjectMeta
</a>
</em>
</td>
<td>
Refer to the Kubernetes API documentation for the fields of the
<code>metadata</code> field.
</td>
</tr>
<tr>
<td>
<code>spec</code></br>
<em>
<a href="#CertificateRequestSpec">
CertificateRequestSpec
</a>
</em>
</td>
<td>
<br/>
<br/>
<table>
<tr>
<td>
<code>commonName</code></br>
<em>
string
</em>
</td>
<td>
<p>Primary domain for the certificate request.</p>
</td>
</tr>
<tr>
<td>
<code>alternateNames</code></br>
<em>
[]string
</em>
</td>
<td>
<p>Additional domains for the certificate request.</p>
</td>
</tr>
</table>
</td>
</tr>
<tr>
<td>
<code>status</code></br>
<em>
<a href="#CertificateRequestStatus">
CertificateRequestStatus
</a>
</em>
</td>
<td>
</td>
</tr>
</tbody>
</table>
<h3 id="CloudFront">CloudFront
</h3>
<p>
<p>CloudFront is the Schema for the cloudfronts API</p>
</p>
<table>
<thead>
<tr>
<th>Field</th>
<th>Description</th>
</tr>
</thead>
<tbody>
<tr>
<td>
<code>apiVersion</code></br>
string</td>
<td>
<code>
aws.skpr.io/v1beta1
</code>
</td>
</tr>
<tr>
<td>
<code>kind</code></br>
string
</td>
<td><code>CloudFront</code></td>
</tr>
<tr>
<td>
<code>metadata</code></br>
<em>
<a href="https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.13/#objectmeta-v1-meta">
Kubernetes meta/v1.ObjectMeta
</a>
</em>
</td>
<td>
Refer to the Kubernetes API documentation for the fields of the
<code>metadata</code> field.
</td>
</tr>
<tr>
<td>
<code>spec</code></br>
<em>
<a href="#CloudFrontSpec">
CloudFrontSpec
</a>
</em>
</td>
<td>
<br/>
<br/>
<table>
<tr>
<td>
<code>aliases</code></br>
<em>
[]string
</em>
</td>
<td>
<p>Aliases which CloudFront will respond to.</p>
</td>
</tr>
<tr>
<td>
<code>certificate</code></br>
<em>
<a href="#CloudFrontSpecCertificate">
CloudFrontSpecCertificate
</a>
</em>
</td>
<td>
<p>Certificate which is applied to this CloudFront distribution.</p>
</td>
</tr>
<tr>
<td>
<code>firewall</code></br>
<em>
<a href="#CloudFrontSpecFirewall">
CloudFrontSpecFirewall
</a>
</em>
</td>
<td>
<p>Firewall configuration for this CloudFront distribution.</p>
</td>
</tr>
<tr>
<td>
<code>behavior</code></br>
<em>
<a href="#CloudFrontSpecBehavior">
CloudFrontSpecBehavior
</a>
</em>
</td>
<td>
<p>Behavior applied to this CloudFront distribution eg. Headers and Cookies.</p>
</td>
</tr>
<tr>
<td>
<code>origin</code></br>
<em>
<a href="#CloudFrontSpecOrigin">
CloudFrontSpecOrigin
</a>
</em>
</td>
<td>
<p>Information CloudFront uses to connect to the backend.</p>
</td>
</tr>
</table>
</td>
</tr>
<tr>
<td>
<code>status</code></br>
<em>
<a href="#CloudFrontStatus">
CloudFrontStatus
</a>
</em>
</td>
<td>
</td>
</tr>
</tbody>
</table>
<h3 id="CloudFrontInvalidation">CloudFrontInvalidation
</h3>
<p>
<p>CloudFrontInvalidation is the Schema for the cloudfrontinvalidations API</p>
</p>
<table>
<thead>
<tr>
<th>Field</th>
<th>Description</th>
</tr>
</thead>
<tbody>
<tr>
<td>
<code>apiVersion</code></br>
string</td>
<td>
<code>
aws.skpr.io/v1beta1
</code>
</td>
</tr>
<tr>
<td>
<code>kind</code></br>
string
</td>
<td><code>CloudFrontInvalidation</code></td>
</tr>
<tr>
<td>
<code>metadata</code></br>
<em>
<a href="https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.13/#objectmeta-v1-meta">
Kubernetes meta/v1.ObjectMeta
</a>
</em>
</td>
<td>
Refer to the Kubernetes API documentation for the fields of the
<code>metadata</code> field.
</td>
</tr>
<tr>
<td>
<code>spec</code></br>
<em>
<a href="#CloudFrontInvalidationSpec">
CloudFrontInvalidationSpec
</a>
</em>
</td>
<td>
<br/>
<br/>
<table>
<tr>
<td>
<code>distribution</code></br>
<em>
string
</em>
</td>
<td>
<p>Name of the CloudFront object.</p>
</td>
</tr>
<tr>
<td>
<code>paths</code></br>
<em>
[]string
</em>
</td>
<td>
<p>Paths which to invalidate.</p>
</td>
</tr>
</table>
</td>
</tr>
<tr>
<td>
<code>status</code></br>
<em>
<a href="#CloudFrontInvalidationStatus">
CloudFrontInvalidationStatus
</a>
</em>
</td>
<td>
</td>
</tr>
</tbody>
</table>
<h3 id="CertificateRequestReference">CertificateRequestReference
</h3>
<p>
(<em>Appears on:</em>
<a href="#CertificateStatus">CertificateStatus</a>)
</p>
<p>
<p>CertificateRequestReference defines the observed state of Certificate</p>
</p>
<table>
<thead>
<tr>
<th>Field</th>
<th>Description</th>
</tr>
</thead>
<tbody>
<tr>
<td>
<code>name</code></br>
<em>
string
</em>
</td>
<td>
<p>Reference name for the certificate request.</p>
</td>
</tr>
<tr>
<td>
<code>details</code></br>
<em>
<a href="#CertificateRequestStatus">
CertificateRequestStatus
</a>
</em>
</td>
<td>
<p>Details of the certificate.</p>
</td>
</tr>
</tbody>
</table>
<h3 id="CertificateRequestSpec">CertificateRequestSpec
</h3>
<p>
(<em>Appears on:</em>
<a href="#CertificateRequest">CertificateRequest</a>, 
<a href="#CertificateSpec">CertificateSpec</a>)
</p>
<p>
<p>CertificateRequestSpec defines the desired state of CertificateRequest</p>
</p>
<table>
<thead>
<tr>
<th>Field</th>
<th>Description</th>
</tr>
</thead>
<tbody>
<tr>
<td>
<code>commonName</code></br>
<em>
string
</em>
</td>
<td>
<p>Primary domain for the certificate request.</p>
</td>
</tr>
<tr>
<td>
<code>alternateNames</code></br>
<em>
[]string
</em>
</td>
<td>
<p>Additional domains for the certificate request.</p>
</td>
</tr>
</tbody>
</table>
<h3 id="CertificateRequestStatus">CertificateRequestStatus
</h3>
<p>
(<em>Appears on:</em>
<a href="#CertificateRequest">CertificateRequest</a>, 
<a href="#CertificateRequestReference">CertificateRequestReference</a>)
</p>
<p>
<p>CertificateRequestStatus defines the observed state of CertificateRequest</p>
</p>
<table>
<thead>
<tr>
<th>Field</th>
<th>Description</th>
</tr>
</thead>
<tbody>
<tr>
<td>
<code>observedGeneration</code></br>
<em>
int64
</em>
</td>
<td>
<p>Used for determining if an APIs information is up to date.</p>
</td>
</tr>
<tr>
<td>
<code>arn</code></br>
<em>
string
</em>
</td>
<td>
<p>Machine identifier for the certificate request.</p>
</td>
</tr>
<tr>
<td>
<code>domains</code></br>
<em>
[]string
</em>
</td>
<td>
<p>Domain list for the certificate.</p>
</td>
</tr>
<tr>
<td>
<code>state</code></br>
<em>
string
</em>
</td>
<td>
<p>Current state of the certificate eg. ISSUED.</p>
</td>
</tr>
<tr>
<td>
<code>validate</code></br>
<em>
<a href="#ValidateRecord">
[]ValidateRecord
</a>
</em>
</td>
<td>
<p>Details used to validate a certificate request.</p>
</td>
</tr>
</tbody>
</table>
<h3 id="CertificateSpec">CertificateSpec
</h3>
<p>
(<em>Appears on:</em>
<a href="#Certificate">Certificate</a>)
</p>
<p>
<p>CertificateSpec defines the desired state of Certificate</p>
</p>
<table>
<thead>
<tr>
<th>Field</th>
<th>Description</th>
</tr>
</thead>
<tbody>
<tr>
<td>
<code>request</code></br>
<em>
<a href="#CertificateRequestSpec">
CertificateRequestSpec
</a>
</em>
</td>
<td>
<p>Information which will be used to provision a certificate.</p>
</td>
</tr>
</tbody>
</table>
<h3 id="CertificateStatus">CertificateStatus
</h3>
<p>
(<em>Appears on:</em>
<a href="#Certificate">Certificate</a>, 
<a href="#IngressStatusCertificateRef">IngressStatusCertificateRef</a>)
</p>
<p>
<p>CertificateStatus defines the observed state of Certificate</p>
</p>
<table>
<thead>
<tr>
<th>Field</th>
<th>Description</th>
</tr>
</thead>
<tbody>
<tr>
<td>
<code>observedGeneration</code></br>
<em>
int64
</em>
</td>
<td>
<p>Used for determining if an APIs information is up to date.</p>
</td>
</tr>
<tr>
<td>
<code>active</code></br>
<em>
<a href="#CertificateRequestReference">
CertificateRequestReference
</a>
</em>
</td>
<td>
<p>The status of the most recently ISSUED certificate.</p>
</td>
</tr>
<tr>
<td>
<code>requests</code></br>
<em>
<a href="#CertificateRequestReference">
[]CertificateRequestReference
</a>
</em>
</td>
<td>
<p>Status of all the certificate requests.</p>
</td>
</tr>
</tbody>
</table>
<h3 id="CloudFrontInvalidationSpec">CloudFrontInvalidationSpec
</h3>
<p>
(<em>Appears on:</em>
<a href="#CloudFrontInvalidation">CloudFrontInvalidation</a>)
</p>
<p>
<p>CloudFrontInvalidationSpec defines the desired state of CloudFrontInvalidation</p>
</p>
<table>
<thead>
<tr>
<th>Field</th>
<th>Description</th>
</tr>
</thead>
<tbody>
<tr>
<td>
<code>distribution</code></br>
<em>
string
</em>
</td>
<td>
<p>Name of the CloudFront object.</p>
</td>
</tr>
<tr>
<td>
<code>paths</code></br>
<em>
[]string
</em>
</td>
<td>
<p>Paths which to invalidate.</p>
</td>
</tr>
</tbody>
</table>
<h3 id="CloudFrontInvalidationStatus">CloudFrontInvalidationStatus
</h3>
<p>
(<em>Appears on:</em>
<a href="#CloudFrontInvalidation">CloudFrontInvalidation</a>)
</p>
<p>
<p>CloudFrontInvalidationStatus defines the observed state of CloudFrontInvalidation</p>
</p>
<table>
<thead>
<tr>
<th>Field</th>
<th>Description</th>
</tr>
</thead>
<tbody>
<tr>
<td>
<code>observedGeneration</code></br>
<em>
int64
</em>
</td>
<td>
<p>Used for determining if an APIs information is up to date.</p>
</td>
</tr>
<tr>
<td>
<code>id</code></br>
<em>
string
</em>
</td>
<td>
<p>Machine identifier for querying an invalidation request.</p>
</td>
</tr>
<tr>
<td>
<code>created</code></br>
<em>
string
</em>
</td>
<td>
<p>When the invalidation request was lodged.</p>
</td>
</tr>
<tr>
<td>
<code>state</code></br>
<em>
string
</em>
</td>
<td>
<p>Current state of the invalidation request.</p>
</td>
</tr>
</tbody>
</table>
<h3 id="CloudFrontSpec">CloudFrontSpec
</h3>
<p>
(<em>Appears on:</em>
<a href="#CloudFront">CloudFront</a>)
</p>
<p>
<p>CloudFrontSpec defines the desired state of CloudFront</p>
</p>
<table>
<thead>
<tr>
<th>Field</th>
<th>Description</th>
</tr>
</thead>
<tbody>
<tr>
<td>
<code>aliases</code></br>
<em>
[]string
</em>
</td>
<td>
<p>Aliases which CloudFront will respond to.</p>
</td>
</tr>
<tr>
<td>
<code>certificate</code></br>
<em>
<a href="#CloudFrontSpecCertificate">
CloudFrontSpecCertificate
</a>
</em>
</td>
<td>
<p>Certificate which is applied to this CloudFront distribution.</p>
</td>
</tr>
<tr>
<td>
<code>firewall</code></br>
<em>
<a href="#CloudFrontSpecFirewall">
CloudFrontSpecFirewall
</a>
</em>
</td>
<td>
<p>Firewall configuration for this CloudFront distribution.</p>
</td>
</tr>
<tr>
<td>
<code>behavior</code></br>
<em>
<a href="#CloudFrontSpecBehavior">
CloudFrontSpecBehavior
</a>
</em>
</td>
<td>
<p>Behavior applied to this CloudFront distribution eg. Headers and Cookies.</p>
</td>
</tr>
<tr>
<td>
<code>origin</code></br>
<em>
<a href="#CloudFrontSpecOrigin">
CloudFrontSpecOrigin
</a>
</em>
</td>
<td>
<p>Information CloudFront uses to connect to the backend.</p>
</td>
</tr>
</tbody>
</table>
<h3 id="CloudFrontSpecBehavior">CloudFrontSpecBehavior
</h3>
<p>
(<em>Appears on:</em>
<a href="#CloudFrontSpec">CloudFrontSpec</a>)
</p>
<p>
<p>CloudFrontSpecBehavior declares the behaviour which will be applied to this CloudFront distribution.</p>
</p>
<table>
<thead>
<tr>
<th>Field</th>
<th>Description</th>
</tr>
</thead>
<tbody>
<tr>
<td>
<code>whitelist</code></br>
<em>
<a href="#CloudFrontSpecBehaviorWhitelist">
CloudFrontSpecBehaviorWhitelist
</a>
</em>
</td>
<td>
<p>Whitelist of headers and cookies.</p>
</td>
</tr>
</tbody>
</table>
<h3 id="CloudFrontSpecBehaviorWhitelist">CloudFrontSpecBehaviorWhitelist
</h3>
<p>
(<em>Appears on:</em>
<a href="#CloudFrontSpecBehavior">CloudFrontSpecBehavior</a>, 
<a href="#IngressSpec">IngressSpec</a>)
</p>
<p>
<p>CloudFrontSpecBehaviorWhitelist declares a whitelist of request parameters which are allowed.</p>
</p>
<table>
<thead>
<tr>
<th>Field</th>
<th>Description</th>
</tr>
</thead>
<tbody>
<tr>
<td>
<code>headers</code></br>
<em>
[]string
</em>
</td>
<td>
<p>Headers which will used when caching.</p>
</td>
</tr>
<tr>
<td>
<code>cookies</code></br>
<em>
[]string
</em>
</td>
<td>
<p>Cookies which will be forwarded to the application.</p>
</td>
</tr>
</tbody>
</table>
<h3 id="CloudFrontSpecCertificate">CloudFrontSpecCertificate
</h3>
<p>
(<em>Appears on:</em>
<a href="#CloudFrontSpec">CloudFrontSpec</a>)
</p>
<p>
<p>CloudFrontSpecCertificate declares a certificate to use for encryption.</p>
</p>
<table>
<thead>
<tr>
<th>Field</th>
<th>Description</th>
</tr>
</thead>
<tbody>
<tr>
<td>
<code>arn</code></br>
<em>
string
</em>
</td>
<td>
<p>Machine identifier for referencing a certificate.</p>
</td>
</tr>
</tbody>
</table>
<h3 id="CloudFrontSpecFirewall">CloudFrontSpecFirewall
</h3>
<p>
(<em>Appears on:</em>
<a href="#CloudFrontSpec">CloudFrontSpec</a>)
</p>
<p>
<p>CloudFrontSpecFirewall declares a firewall which this CloudFront is associated with.</p>
</p>
<table>
<thead>
<tr>
<th>Field</th>
<th>Description</th>
</tr>
</thead>
<tbody>
<tr>
<td>
<code>arn</code></br>
<em>
string
</em>
</td>
<td>
<p>Machine identifier for referencing a firewall.</p>
</td>
</tr>
</tbody>
</table>
<h3 id="CloudFrontSpecOrigin">CloudFrontSpecOrigin
</h3>
<p>
(<em>Appears on:</em>
<a href="#CloudFrontSpec">CloudFrontSpec</a>)
</p>
<p>
<p>CloudFrontSpecOrigin declares the origin which traffic will be sent.</p>
</p>
<table>
<thead>
<tr>
<th>Field</th>
<th>Description</th>
</tr>
</thead>
<tbody>
<tr>
<td>
<code>endpoint</code></br>
<em>
string
</em>
</td>
<td>
<p>Backend connection information for CloudFront.</p>
</td>
</tr>
<tr>
<td>
<code>policy</code></br>
<em>
string
</em>
</td>
<td>
<p>Backend connection information for CloudFront.</p>
</td>
</tr>
<tr>
<td>
<code>timeout</code></br>
<em>
int64
</em>
</td>
<td>
<p>How long CloudFront should wait before timing out.</p>
</td>
</tr>
</tbody>
</table>
<h3 id="CloudFrontStatus">CloudFrontStatus
</h3>
<p>
(<em>Appears on:</em>
<a href="#CloudFront">CloudFront</a>, 
<a href="#IngressStatusCloudFrontRef">IngressStatusCloudFrontRef</a>)
</p>
<p>
<p>CloudFrontStatus defines the observed state of CloudFront</p>
</p>
<table>
<thead>
<tr>
<th>Field</th>
<th>Description</th>
</tr>
</thead>
<tbody>
<tr>
<td>
<code>observedGeneration</code></br>
<em>
int64
</em>
</td>
<td>
<p>Used for determining if an APIs information is up to date.</p>
</td>
</tr>
<tr>
<td>
<code>arn</code></br>
<em>
string
</em>
</td>
<td>
<p>Machine identifier for querying the CloudFront distribution.</p>
</td>
</tr>
<tr>
<td>
<code>state</code></br>
<em>
string
</em>
</td>
<td>
<p>Current state of the CloudFront distribution.</p>
</td>
</tr>
<tr>
<td>
<code>domainName</code></br>
<em>
string
</em>
</td>
<td>
<p>DomainName for creating CNAME records.</p>
</td>
</tr>
<tr>
<td>
<code>runningInvalidations</code></br>
<em>
int64
</em>
</td>
<td>
<p>How many invalidations are currently running.</p>
</td>
</tr>
</tbody>
</table>
<h3 id="ValidateRecord">ValidateRecord
</h3>
<p>
(<em>Appears on:</em>
<a href="#CertificateRequestStatus">CertificateRequestStatus</a>)
</p>
<p>
<p>ValidateRecord provide details to site administrators on how to validate a certificate.</p>
</p>
<table>
<thead>
<tr>
<th>Field</th>
<th>Description</th>
</tr>
</thead>
<tbody>
<tr>
<td>
<code>arn</code></br>
<em>
string
</em>
</td>
<td>
<p>The name of DNS validation record.</p>
</td>
</tr>
<tr>
<td>
<code>type</code></br>
<em>
string
</em>
</td>
<td>
<p>The type of DNS validation record.</p>
</td>
</tr>
<tr>
<td>
<code>value</code></br>
<em>
string
</em>
</td>
<td>
<p>The value of DNS validation record.</p>
</td>
</tr>
</tbody>
</table>
<hr/>
<h2 id="edge.skpr.io">edge.skpr.io</h2>
<p>
<p>Package v1beta1 contains API Schema definitions for the edge v1beta1 API group</p>
</p>
Resource Types:
<ul><li>
<a href="#Ingress">Ingress</a>
</li></ul>
<h3 id="Ingress">Ingress
</h3>
<p>
<p>Ingress is the Schema for the ingresses API</p>
</p>
<table>
<thead>
<tr>
<th>Field</th>
<th>Description</th>
</tr>
</thead>
<tbody>
<tr>
<td>
<code>apiVersion</code></br>
string</td>
<td>
<code>
edge.skpr.io/v1beta1
</code>
</td>
</tr>
<tr>
<td>
<code>kind</code></br>
string
</td>
<td><code>Ingress</code></td>
</tr>
<tr>
<td>
<code>metadata</code></br>
<em>
<a href="https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.13/#objectmeta-v1-meta">
Kubernetes meta/v1.ObjectMeta
</a>
</em>
</td>
<td>
Refer to the Kubernetes API documentation for the fields of the
<code>metadata</code> field.
</td>
</tr>
<tr>
<td>
<code>spec</code></br>
<em>
<a href="#IngressSpec">
IngressSpec
</a>
</em>
</td>
<td>
<br/>
<br/>
<table>
<tr>
<td>
<code>routes</code></br>
<em>
<a href="#IngressSpecRoutes">
IngressSpecRoutes
</a>
</em>
</td>
<td>
<p>Rules which are used to Ingress traffic to an application.</p>
</td>
</tr>
<tr>
<td>
<code>whitelist</code></br>
<em>
<a href="#CloudFrontSpecBehaviorWhitelist">
CloudFrontSpecBehaviorWhitelist
</a>
</em>
</td>
<td>
<p>Whitelist rules for CloudFront.</p>
</td>
</tr>
<tr>
<td>
<code>service</code></br>
<em>
<a href="#IngressSpecService">
IngressSpecService
</a>
</em>
</td>
<td>
<p>Backend connectivity details.</p>
</td>
</tr>
</table>
</td>
</tr>
<tr>
<td>
<code>status</code></br>
<em>
<a href="#IngressStatus">
IngressStatus
</a>
</em>
</td>
<td>
</td>
</tr>
</tbody>
</table>
<h3 id="IngressSpec">IngressSpec
</h3>
<p>
(<em>Appears on:</em>
<a href="#Ingress">Ingress</a>)
</p>
<p>
<p>IngressSpec defines the desired state of Ingress</p>
</p>
<table>
<thead>
<tr>
<th>Field</th>
<th>Description</th>
</tr>
</thead>
<tbody>
<tr>
<td>
<code>routes</code></br>
<em>
<a href="#IngressSpecRoutes">
IngressSpecRoutes
</a>
</em>
</td>
<td>
<p>Rules which are used to Ingress traffic to an application.</p>
</td>
</tr>
<tr>
<td>
<code>whitelist</code></br>
<em>
<a href="#CloudFrontSpecBehaviorWhitelist">
CloudFrontSpecBehaviorWhitelist
</a>
</em>
</td>
<td>
<p>Whitelist rules for CloudFront.</p>
</td>
</tr>
<tr>
<td>
<code>service</code></br>
<em>
<a href="#IngressSpecService">
IngressSpecService
</a>
</em>
</td>
<td>
<p>Backend connectivity details.</p>
</td>
</tr>
</tbody>
</table>
<h3 id="IngressSpecRoute">IngressSpecRoute
</h3>
<p>
(<em>Appears on:</em>
<a href="#IngressSpecRoutes">IngressSpecRoutes</a>)
</p>
<p>
<p>IngressSpecRoute traffic from a domain and path to a service.</p>
</p>
<table>
<thead>
<tr>
<th>Field</th>
<th>Description</th>
</tr>
</thead>
<tbody>
<tr>
<td>
<code>domain</code></br>
<em>
string
</em>
</td>
<td>
<p>Domain used as part of a route rule.</p>
</td>
</tr>
<tr>
<td>
<code>subpaths</code></br>
<em>
[]string
</em>
</td>
<td>
<p>Supaths included in the route rule.</p>
</td>
</tr>
</tbody>
</table>
<h3 id="IngressSpecRoutes">IngressSpecRoutes
</h3>
<p>
(<em>Appears on:</em>
<a href="#IngressSpec">IngressSpec</a>)
</p>
<p>
<p>IngressSpecRoutes declare the routes for the application.</p>
</p>
<table>
<thead>
<tr>
<th>Field</th>
<th>Description</th>
</tr>
</thead>
<tbody>
<tr>
<td>
<code>primary</code></br>
<em>
<a href="#IngressSpecRoute">
IngressSpecRoute
</a>
</em>
</td>
<td>
<p>Primary domain and routing rule for the application.</p>
</td>
</tr>
<tr>
<td>
<code>secondary</code></br>
<em>
<a href="#IngressSpecRoute">
[]IngressSpecRoute
</a>
</em>
</td>
<td>
<p>Seconard domains and routing rules for the application.</p>
</td>
</tr>
</tbody>
</table>
<h3 id="IngressSpecService">IngressSpecService
</h3>
<p>
(<em>Appears on:</em>
<a href="#IngressSpec">IngressSpec</a>)
</p>
<p>
<p>IngressSpecService connects an Ingress to a Service.</p>
</p>
<table>
<thead>
<tr>
<th>Field</th>
<th>Description</th>
</tr>
</thead>
<tbody>
<tr>
<td>
<code>name</code></br>
<em>
string
</em>
</td>
<td>
<p>Name of the Kubernetes Service object to route traffic to.</p>
</td>
</tr>
<tr>
<td>
<code>port</code></br>
<em>
int
</em>
</td>
<td>
<p>Port of the Kubernetes Service object to route traffic to.</p>
</td>
</tr>
</tbody>
</table>
<h3 id="IngressStatus">IngressStatus
</h3>
<p>
(<em>Appears on:</em>
<a href="#Ingress">Ingress</a>)
</p>
<p>
<p>IngressStatus defines the observed state of Ingress</p>
</p>
<table>
<thead>
<tr>
<th>Field</th>
<th>Description</th>
</tr>
</thead>
<tbody>
<tr>
<td>
<code>observedGeneration</code></br>
<em>
int64
</em>
</td>
<td>
<p>Used for determining if an APIs information is up to date.</p>
</td>
</tr>
<tr>
<td>
<code>CloudFront</code></br>
<em>
<a href="#IngressStatusCloudFrontRef">
IngressStatusCloudFrontRef
</a>
</em>
</td>
<td>
<p>Status of the provisioned CloudFront distribution.</p>
</td>
</tr>
<tr>
<td>
<code>Certificate</code></br>
<em>
<a href="#IngressStatusCertificateRef">
IngressStatusCertificateRef
</a>
</em>
</td>
<td>
<p>Status of the provisioned Certificate.</p>
</td>
</tr>
</tbody>
</table>
<h3 id="IngressStatusCertificateRef">IngressStatusCertificateRef
</h3>
<p>
(<em>Appears on:</em>
<a href="#IngressStatus">IngressStatus</a>)
</p>
<p>
<p>IngressStatusCertificateRef provides status on the provisioned Certificate.</p>
</p>
<table>
<thead>
<tr>
<th>Field</th>
<th>Description</th>
</tr>
</thead>
<tbody>
<tr>
<td>
<code>Name</code></br>
<em>
string
</em>
</td>
<td>
<p>Name of the certificate.</p>
</td>
</tr>
<tr>
<td>
<code>Details</code></br>
<em>
<a href="#CertificateStatus">
CertificateStatus
</a>
</em>
</td>
<td>
<p>Details on the provisioned certificate.</p>
</td>
</tr>
</tbody>
</table>
<h3 id="IngressStatusCloudFrontRef">IngressStatusCloudFrontRef
</h3>
<p>
(<em>Appears on:</em>
<a href="#IngressStatus">IngressStatus</a>)
</p>
<p>
<p>IngressStatusCloudFrontRef provides status on the provisioned CloudFront.</p>
</p>
<table>
<thead>
<tr>
<th>Field</th>
<th>Description</th>
</tr>
</thead>
<tbody>
<tr>
<td>
<code>name</code></br>
<em>
string
</em>
</td>
<td>
<p>Name of the CloudFront distribution.</p>
</td>
</tr>
<tr>
<td>
<code>Details</code></br>
<em>
<a href="#CloudFrontStatus">
CloudFrontStatus
</a>
</em>
</td>
<td>
<p>Details on the provisioned CloudFront distribution.</p>
</td>
</tr>
</tbody>
</table>
<hr/>
<h2 id="extensions.skpr.io">extensions.skpr.io</h2>
<p>
<p>Package v1beta1 contains API Schema definitions for the extensions v1beta1 API group</p>
</p>
Resource Types:
<ul><li>
<a href="#Exec">Exec</a>
</li><li>
<a href="#SMTP">SMTP</a>
</li></ul>
<h3 id="Exec">Exec
</h3>
<p>
<p>Exec is the Schema for the execs API</p>
</p>
<table>
<thead>
<tr>
<th>Field</th>
<th>Description</th>
</tr>
</thead>
<tbody>
<tr>
<td>
<code>apiVersion</code></br>
string</td>
<td>
<code>
extensions.skpr.io/v1beta1
</code>
</td>
</tr>
<tr>
<td>
<code>kind</code></br>
string
</td>
<td><code>Exec</code></td>
</tr>
<tr>
<td>
<code>metadata</code></br>
<em>
<a href="https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.13/#objectmeta-v1-meta">
Kubernetes meta/v1.ObjectMeta
</a>
</em>
</td>
<td>
Refer to the Kubernetes API documentation for the fields of the
<code>metadata</code> field.
</td>
</tr>
<tr>
<td>
<code>spec</code></br>
<em>
<a href="#ExecSpec">
ExecSpec
</a>
</em>
</td>
<td>
<br/>
<br/>
<table>
<tr>
<td>
<code>entrypoint</code></br>
<em>
string
</em>
</td>
<td>
<p>Container which commands will be executed.</p>
</td>
</tr>
<tr>
<td>
<code>template</code></br>
<em>
<a href="https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.13/#podspec-v1-core">
Kubernetes core/v1.PodSpec
</a>
</em>
</td>
<td>
<p>Template used when provisioning an execution environment.</p>
</td>
</tr>
</table>
</td>
</tr>
</tbody>
</table>
<h3 id="SMTP">SMTP
</h3>
<p>
<p>SMTP is the Schema for the smtps API</p>
</p>
<table>
<thead>
<tr>
<th>Field</th>
<th>Description</th>
</tr>
</thead>
<tbody>
<tr>
<td>
<code>apiVersion</code></br>
string</td>
<td>
<code>
extensions.skpr.io/v1beta1
</code>
</td>
</tr>
<tr>
<td>
<code>kind</code></br>
string
</td>
<td><code>SMTP</code></td>
</tr>
<tr>
<td>
<code>metadata</code></br>
<em>
<a href="https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.13/#objectmeta-v1-meta">
Kubernetes meta/v1.ObjectMeta
</a>
</em>
</td>
<td>
Refer to the Kubernetes API documentation for the fields of the
<code>metadata</code> field.
</td>
</tr>
<tr>
<td>
<code>spec</code></br>
<em>
<a href="#SMTPSpec">
SMTPSpec
</a>
</em>
</td>
<td>
<br/>
<br/>
<table>
<tr>
<td>
<code>from</code></br>
<em>
<a href="#SMTPSpecFrom">
SMTPSpecFrom
</a>
</em>
</td>
<td>
<p>From defines what an application is allowed to send from.</p>
</td>
</tr>
</table>
</td>
</tr>
<tr>
<td>
<code>status</code></br>
<em>
<a href="#SMTPStatus">
SMTPStatus
</a>
</em>
</td>
<td>
</td>
</tr>
</tbody>
</table>
<h3 id="ExecSpec">ExecSpec
</h3>
<p>
(<em>Appears on:</em>
<a href="#Exec">Exec</a>)
</p>
<p>
<p>ExecSpec defines the desired state of Exec</p>
</p>
<table>
<thead>
<tr>
<th>Field</th>
<th>Description</th>
</tr>
</thead>
<tbody>
<tr>
<td>
<code>entrypoint</code></br>
<em>
string
</em>
</td>
<td>
<p>Container which commands will be executed.</p>
</td>
</tr>
<tr>
<td>
<code>template</code></br>
<em>
<a href="https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.13/#podspec-v1-core">
Kubernetes core/v1.PodSpec
</a>
</em>
</td>
<td>
<p>Template used when provisioning an execution environment.</p>
</td>
</tr>
</tbody>
</table>
<h3 id="SMTPSpec">SMTPSpec
</h3>
<p>
(<em>Appears on:</em>
<a href="#SMTP">SMTP</a>)
</p>
<p>
<p>SMTPSpec defines the desired state of SMTP</p>
</p>
<table>
<thead>
<tr>
<th>Field</th>
<th>Description</th>
</tr>
</thead>
<tbody>
<tr>
<td>
<code>from</code></br>
<em>
<a href="#SMTPSpecFrom">
SMTPSpecFrom
</a>
</em>
</td>
<td>
<p>From defines what an application is allowed to send from.</p>
</td>
</tr>
</tbody>
</table>
<h3 id="SMTPSpecFrom">SMTPSpecFrom
</h3>
<p>
(<em>Appears on:</em>
<a href="#DrupalSpecSMTP">DrupalSpecSMTP</a>, 
<a href="#SMTPSpec">SMTPSpec</a>)
</p>
<p>
<p>SMTPSpecFrom defines what an application is allowed to send from.</p>
</p>
<table>
<thead>
<tr>
<th>Field</th>
<th>Description</th>
</tr>
</thead>
<tbody>
<tr>
<td>
<code>address</code></br>
<em>
string
</em>
</td>
<td>
<p>Address which an application is allowed to send from.</p>
</td>
</tr>
</tbody>
</table>
<h3 id="SMTPStatus">SMTPStatus
</h3>
<p>
(<em>Appears on:</em>
<a href="#SMTP">SMTP</a>)
</p>
<p>
<p>SMTPStatus defines the observed state of SMTP</p>
</p>
<table>
<thead>
<tr>
<th>Field</th>
<th>Description</th>
</tr>
</thead>
<tbody>
<tr>
<td>
<code>verification</code></br>
<em>
<a href="#SMTPStatusVerification">
SMTPStatusVerification
</a>
</em>
</td>
<td>
<p>Provides the status of verifying FROM attributes.</p>
</td>
</tr>
<tr>
<td>
<code>connection</code></br>
<em>
<a href="#SMTPStatusConnection">
SMTPStatusConnection
</a>
</em>
</td>
<td>
<p>Provides connection details for sending email.</p>
</td>
</tr>
</tbody>
</table>
<h3 id="SMTPStatusConnection">SMTPStatusConnection
</h3>
<p>
(<em>Appears on:</em>
<a href="#SMTPStatus">SMTPStatus</a>)
</p>
<p>
<p>SMTPStatusConnection provides connection details for sending email.</p>
</p>
<table>
<thead>
<tr>
<th>Field</th>
<th>Description</th>
</tr>
</thead>
<tbody>
<tr>
<td>
<code>hostname</code></br>
<em>
string
</em>
</td>
<td>
<p>Hostname used when connecting to the SMTP server.</p>
</td>
</tr>
<tr>
<td>
<code>port</code></br>
<em>
int
</em>
</td>
<td>
<p>Port used when connecting to the SMTP server.</p>
</td>
</tr>
<tr>
<td>
<code>username</code></br>
<em>
string
</em>
</td>
<td>
<p>Username used when connecting to the SMTP server.</p>
</td>
</tr>
<tr>
<td>
<code>password</code></br>
<em>
string
</em>
</td>
<td>
<p>Password used when connecting to the SMTP server.</p>
</td>
</tr>
</tbody>
</table>
<h3 id="SMTPStatusVerification">SMTPStatusVerification
</h3>
<p>
(<em>Appears on:</em>
<a href="#DrupalStatusSMTP">DrupalStatusSMTP</a>, 
<a href="#SMTPStatus">SMTPStatus</a>)
</p>
<p>
<p>SMTPStatusVerification provides the status of verifying FROM attributes.</p>
</p>
<table>
<thead>
<tr>
<th>Field</th>
<th>Description</th>
</tr>
</thead>
<tbody>
<tr>
<td>
<code>address</code></br>
<em>
string
</em>
</td>
<td>
<p>Address which an application is allowed to send from.</p>
</td>
</tr>
</tbody>
</table>
<hr/>
<h2 id="mysql.skpr.io">mysql.skpr.io</h2>
<p>
<p>Package v1beta1 contains API Schema definitions for the mysql v1beta1 API group</p>
</p>
Resource Types:
<ul><li>
<a href="#Database">Database</a>
</li></ul>
<h3 id="Database">Database
</h3>
<p>
<p>Database is the Schema for the databaseclaims API</p>
</p>
<table>
<thead>
<tr>
<th>Field</th>
<th>Description</th>
</tr>
</thead>
<tbody>
<tr>
<td>
<code>apiVersion</code></br>
string</td>
<td>
<code>
mysql.skpr.io/v1beta1
</code>
</td>
</tr>
<tr>
<td>
<code>kind</code></br>
string
</td>
<td><code>Database</code></td>
</tr>
<tr>
<td>
<code>metadata</code></br>
<em>
<a href="https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.13/#objectmeta-v1-meta">
Kubernetes meta/v1.ObjectMeta
</a>
</em>
</td>
<td>
Refer to the Kubernetes API documentation for the fields of the
<code>metadata</code> field.
</td>
</tr>
<tr>
<td>
<code>spec</code></br>
<em>
<a href="#DatabaseSpec">
DatabaseSpec
</a>
</em>
</td>
<td>
<br/>
<br/>
<table>
<tr>
<td>
<code>provisioner</code></br>
<em>
string
</em>
</td>
<td>
<p>Provisioner used to create databases.</p>
</td>
</tr>
<tr>
<td>
<code>privileges</code></br>
<em>
[]string
</em>
</td>
<td>
<p>Privileges which the application requires.</p>
</td>
</tr>
</table>
</td>
</tr>
<tr>
<td>
<code>status</code></br>
<em>
<a href="#DatabaseStatus">
DatabaseStatus
</a>
</em>
</td>
<td>
</td>
</tr>
</tbody>
</table>
<h3 id="DatabaseSpec">DatabaseSpec
</h3>
<p>
(<em>Appears on:</em>
<a href="#Database">Database</a>)
</p>
<p>
<p>DatabaseSpec defines the desired state of Database</p>
</p>
<table>
<thead>
<tr>
<th>Field</th>
<th>Description</th>
</tr>
</thead>
<tbody>
<tr>
<td>
<code>provisioner</code></br>
<em>
string
</em>
</td>
<td>
<p>Provisioner used to create databases.</p>
</td>
</tr>
<tr>
<td>
<code>privileges</code></br>
<em>
[]string
</em>
</td>
<td>
<p>Privileges which the application requires.</p>
</td>
</tr>
</tbody>
</table>
<h3 id="DatabaseStatus">DatabaseStatus
</h3>
<p>
(<em>Appears on:</em>
<a href="#Database">Database</a>)
</p>
<p>
<p>DatabaseStatus defines the observed state of Database</p>
</p>
<table>
<thead>
<tr>
<th>Field</th>
<th>Description</th>
</tr>
</thead>
<tbody>
<tr>
<td>
<code>observedGeneration</code></br>
<em>
int64
</em>
</td>
<td>
<p>Used for determining if an APIs information is up to date.</p>
</td>
</tr>
<tr>
<td>
<code>phase</code></br>
<em>
<a href="#Phase">
Phase
</a>
</em>
</td>
<td>
<p>Current state of the database being provisioned.</p>
</td>
</tr>
<tr>
<td>
<code>connection</code></br>
<em>
<a href="#DatabaseStatusConnection">
DatabaseStatusConnection
</a>
</em>
</td>
<td>
<p>Connection details for the database.</p>
</td>
</tr>
</tbody>
</table>
<h3 id="DatabaseStatusConnection">DatabaseStatusConnection
</h3>
<p>
(<em>Appears on:</em>
<a href="#DatabaseStatus">DatabaseStatus</a>)
</p>
<p>
<p>DatabaseStatusConnection for applications.</p>
</p>
<table>
<thead>
<tr>
<th>Field</th>
<th>Description</th>
</tr>
</thead>
<tbody>
<tr>
<td>
<code>hostname</code></br>
<em>
string
</em>
</td>
<td>
<p>Hostname used when connecting to the database.</p>
</td>
</tr>
<tr>
<td>
<code>port</code></br>
<em>
int
</em>
</td>
<td>
<p>Port used when connecting to the database.</p>
</td>
</tr>
<tr>
<td>
<code>database</code></br>
<em>
string
</em>
</td>
<td>
<p>Database used when connecting to the database.</p>
</td>
</tr>
<tr>
<td>
<code>username</code></br>
<em>
string
</em>
</td>
<td>
<p>Username used when connecting to the database.</p>
</td>
</tr>
<tr>
<td>
<code>password</code></br>
<em>
string
</em>
</td>
<td>
<p>Password used when connecting to the database.</p>
</td>
</tr>
</tbody>
</table>
<h3 id="Phase">Phase
(<code>string</code> alias)</p></h3>
<p>
(<em>Appears on:</em>
<a href="#DatabaseStatus">DatabaseStatus</a>)
</p>
<p>
<p>Phase which indicates the status of an object.</p>
</p>
<hr/>
<p><em>
Generated with <code>gen-crd-api-reference-docs</code>
on git commit <code>76309af</code>.
</em></p>
