package cmos.middleware.rsp.controller;

import org.springframework.web.bind.annotation.RestController;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;

import io.kubernetes.client.custom.V1Patch;
import io.kubernetes.client.openapi.ApiException;
import io.kubernetes.client.openapi.models.V1Deployment;
import io.kubernetes.client.openapi.models.V1DeploymentList;
import io.kubernetes.client.openapi.models.V1Status;
import io.kubernetes.client.util.PatchUtils;
import io.swagger.annotations.*;
import io.kubernetes.client.openapi.ApiResponse;

import com.alibaba.fastjson.JSON;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.validation.annotation.Validated;
import org.springframework.web.bind.annotation.*;

import java.util.ArrayList;
import java.util.HashMap;
import java.util.List;
import java.util.Map;

@RestController
@RequestMapping(value = "/")
@Api(tags = "Controller", description = "Deployment管理")
@Validated
public class Controller {
	@Autowired
	KubeClient kubeclient;

	String namespace = "panji-redis";

	private static final Logger logger = LoggerFactory.getLogger(Controller.class);

	@RequestMapping(value = "/deployments", method = RequestMethod.GET)
	@ApiOperation(value = "获取集群列表", notes = "获取集群列表")
	@ResponseBody
	public ResponseData<List<V1Deployment>> list() {
		ResponseData<List<V1Deployment>> resp = new ResponseData<List<V1Deployment>>();

		try {
			V1DeploymentList list = kubeclient.getAppsV1Api().listNamespacedDeployment(
				namespace, null, null, null, null, null, 100, null, null, 10, null
			);
			resp.setData(list.getItems());
		} catch (ApiException e) {
			resp.setStatus(1);
			resp.setMessage(e.getMessage());
		}
		return resp;
	}

	@RequestMapping(value = "/deployments/{name}", method = RequestMethod.GET)
	@ApiOperation(value = "获取集群详情", notes = "获取集群详情")
	@ResponseBody
	public ResponseData<V1Deployment> detail(@PathVariable String name) {
		ResponseData<V1Deployment> resp = new ResponseData<V1Deployment>();
		try {
			V1Deployment deploy = kubeclient.getAppsV1Api().readNamespacedDeployment(
				name, namespace, null
			);
			resp.setData(deploy);
		} catch (ApiException e) {
			resp.setStatus(1);
			resp.setMessage(e.getMessage());
		}
		return resp;
	}

	@RequestMapping(value = "/deployments/{name}", method = RequestMethod.PUT)
	@ApiOperation(value = "集群变更", notes = "集群变更")
	@ResponseBody
	public ResponseData<String> update(
		@PathVariable String name,
		@RequestBody V1Deployment newDeploy
	) {
		ResponseData<String> resp = new ResponseData<String>();

		// String newDeployStr = Configuration.getDefaultApiClient().getJSON().serialize(newDeploy);
		String newDeployStr = kubeclient.getKubeClient().getJSON().serialize(newDeploy);
		try {
			// 这种形式只能合并(新增和变更), 对已有 labels 不能通过设置为 null 进行清空.
			// (指定为 PATCH_FORMAT_APPLY_YAML 模式也不行)
			V1Patch patch = new V1Patch(newDeployStr);
			PatchUtils.patch(
				Object.class,
				() -> kubeclient.getAppsV1Api().patchNamespacedDeploymentCall(
					name, namespace, patch, 
					// force 字段, 不要设置为 false, 必须指定为 null, 否则可能出现 "Unprocessable Entity"
					null, null, null, null, null, null
				),
				V1Patch.PATCH_FORMAT_JSON_MERGE_PATCH,
				kubeclient.getKubeClient()
			);
			logger.info("new deployment: {}", newDeploy);
		} catch (ApiException e) {
			resp.setStatus(1);
			resp.setMessage(e.getMessage());
		}
		return resp;
	}

	@RequestMapping(value = "/deployments_2/{name}", method = RequestMethod.PUT)
	@ApiOperation(value = "集群变更2", notes = "集群变更2")
	@ResponseBody
	public ResponseData<String> update2(
		@PathVariable String name, @RequestBody V1Deployment newDeploy
	) {
		ResponseData<String> resp = new ResponseData<String>();

        List<Object> patchBody = new ArrayList<>();
		// 尝试通过 path 设置为 /, value 为 newDeploy 对整个 deployment 进行更新, 无效.
		// 不过分段进行更新则是可以的, 可以实现类似 kubectl apply 的效果, 
		// 对 labels, annotations 这种 map 类似也能实现清空的效果.
        Map<String, Object> metaOperator = new HashMap<>();
		// 可用的 op 有 [add, replace, remove]
        metaOperator.put("op", "replace");
		metaOperator.put("path", "/metadata");
		metaOperator.put("value", newDeploy.getMetadata());
        patchBody.add(metaOperator);
		Map<String, Object> specOperator = new HashMap<>();
        specOperator.put("op", "replace");
		specOperator.put("path", "/spec");
		specOperator.put("value", newDeploy.getSpec());
        patchBody.add(specOperator);
		V1Patch patch = new V1Patch(JSON.toJSONString(patchBody));

		try {
			newDeploy = kubeclient.getAppsV1Api().patchNamespacedDeployment(
				name, namespace, patch,
				// force 字段, 不要设置为 false, 必须指定为 null, 否则可能出现 "Unprocessable Entity"
				null, null, null, null, null 
			);
			logger.info("new deployment: {}", newDeploy);
		} catch (ApiException e) {
			resp.setStatus(1);
			resp.setMessage(e.getMessage());
		}

		return resp;
	}

	// 集群变更(废弃)
	public ResponseData<String> update3(@PathVariable String name, @RequestBody V1Deployment newDeploy) {
		ResponseData<String> resp = new ResponseData<String>();

		// 这种更新方式会报错(patch 对象不能这样构建)
		// json: cannot unmarshal object into Go value of type jsonpatch.Patch","reason":"BadRequest","code":400"

		String newDeployStr = kubeclient.getKubeClient().getJSON().serialize(newDeploy);
		V1Patch patch = new V1Patch(newDeployStr);
		try {
			// force 字段, 不要设置为 false, 必须指定为 null, 否则可能出现 "Unprocessable Entity"
			V1Deployment deploy1 = kubeclient.getAppsV1Api().patchNamespacedDeployment(
				name, namespace, patch, null, null, null, null, null
			);
			ApiResponse<V1Deployment> deploy2 = kubeclient.getAppsV1Api().patchNamespacedDeploymentWithHttpInfo(
				name, namespace, patch, null, null, null, null, null
			);
			logger.info("new deployment: {}", deploy1);
			logger.info("new deployment: {}", deploy2);
		} catch (ApiException e) {
			resp.setStatus(1);
			resp.setMessage(e.getMessage());
		}
		return resp;
	}

	@RequestMapping(value = "/deployments/{name}/labels1", method = RequestMethod.PUT)
	@ApiOperation(value = "变更labels", notes = "变更labels")
	@ResponseBody
	public ResponseData<String> labels1(
		@PathVariable String name, @RequestBody Map<String, String> newLabels
	) {
		ResponseData<String> resp = new ResponseData<String>();

        List<Object> patchBody = new ArrayList<>();
		// 可用的 op 有 [add, replace, remove]
		for (Map.Entry<String, String> entry : newLabels.entrySet()) {
			Map<String, String> labelOperator = new HashMap<>();
			labelOperator.put("op", "add");
			labelOperator.put("path", "/metadata/labels/" + entry.getKey());
			labelOperator.put("value", entry.getValue());
			patchBody.add(labelOperator);
		}
		V1Patch patch = new V1Patch(JSON.toJSONString(patchBody));

		try {
			V1Deployment newDeploy = kubeclient.getAppsV1Api().patchNamespacedDeployment(
				name, namespace, patch,
				// force 字段, 不要设置为 false, 必须指定为 null, 否则可能出现 "Unprocessable Entity"
				null, null, null, null, null 
			);
			logger.info("new deployment: {}", newDeploy);
		} catch (ApiException e) {
			resp.setStatus(1);
			resp.setMessage(e.getMessage());
		}
		return resp;
	}

	@RequestMapping(value = "/deployments/{name}/labels2", method = RequestMethod.PUT)
	@ApiOperation(value = "替换labels", notes = "替换labels")
	@ResponseBody
	public ResponseData<String> labels2(
		@PathVariable String name, @RequestBody Map<String, String> newLabels
	) {
		ResponseData<String> resp = new ResponseData<String>();

        List<Object> patchBody = new ArrayList<>();
        Map<String, Object> labelOperator = new HashMap<>();
		// 可用的 op 有 [add, replace, remove]
		//
		// 尝试通过 path 设置为 /, value 为 newDeploy 对整个 deployment 进行更新, 无效.
        labelOperator.put("op", "replace");
		labelOperator.put("path", "/metadata/labels");
		labelOperator.put("value", newLabels);
        patchBody.add(labelOperator);
		V1Patch patch = new V1Patch(JSON.toJSONString(patchBody));

		try {
			V1Deployment newDeploy = kubeclient.getAppsV1Api().patchNamespacedDeployment(
				name, namespace, patch,
				// force 字段, 不要设置为 false, 必须指定为 null, 否则可能出现 "Unprocessable Entity"
				null, null, null, null, null 
			);
			logger.info("new deployment: {}", newDeploy);
		} catch (ApiException e) {
			resp.setStatus(1);
			resp.setMessage(e.getMessage());
		}
		return resp;
	}

	@RequestMapping(value = "/deployments/{name}", method = RequestMethod.DELETE)
	@ApiOperation(value = "删除集群", notes = "删除集群")
	@ResponseBody
	public ResponseData<V1Status> delete(@PathVariable String name) {
		ResponseData<V1Status> resp = new ResponseData<V1Status>();

		try {
			V1Status status = kubeclient.getAppsV1Api().deleteNamespacedDeployment(
				name,
				namespace,
				null, // String pretty, 
				null, // String dryRun,
				0, // Integer gracePeriodSeconds,
				false, // Boolean orphanDependents,
				null, // String propagationPolicy,
				null // V1DeleteOptions body
			);
			resp.setData(status);
			logger.info("delete result: {}", status);
		} catch (ApiException e) {
			resp.setStatus(1);
			resp.setMessage(e.getMessage());
		}
		return resp;
	}
}
