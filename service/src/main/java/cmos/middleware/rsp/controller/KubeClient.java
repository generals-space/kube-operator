package cmos.middleware.rsp.controller;

import java.util.concurrent.TimeUnit;

import org.springframework.stereotype.Component;

import io.kubernetes.client.openapi.ApiClient;
import io.kubernetes.client.openapi.apis.AppsV1Api;
import io.kubernetes.client.openapi.apis.CoreV1Api;
import io.kubernetes.client.openapi.apis.CustomObjectsApi;
import io.kubernetes.client.util.ClientBuilder;
import io.kubernetes.client.util.Config;
import io.kubernetes.client.util.credentials.AccessTokenAuthentication;
import okhttp3.OkHttpClient;

@Component
public class KubeClient {
    String apiserverAddr = "https://192.168.29.104:6443";
    String apiserverToken = "";

    public ApiClient getKubeClient() {
        ApiClient client = Config.fromToken(apiserverAddr, apiserverToken, false);
        OkHttpClient httpClient = client.getHttpClient().newBuilder().readTimeout(0, TimeUnit.SECONDS).build();
        client.setHttpClient(httpClient);
        // client.setDebugging(true);
        return client;
    }

    public ApiClient getKubeClient2() {
        // 以下3个参数与 getKubeClient() 中 Config.fromToken() 的3个参数是一致的.
        ApiClient client = new ClientBuilder()
            //设置 k8s 服务所在 ip地址
            .setBasePath(apiserverAddr)
            //是否开启 ssl 验证
            .setVerifyingSsl(false)
            //插入访问连接用的 Token
            .setAuthentication(new AccessTokenAuthentication(apiserverToken))
            .build();
        io.kubernetes.client.openapi.Configuration.setDefaultApiClient(client);
        // client.setDebugging(true);
        return client;
    }

    public CustomObjectsApi getCrdApi() {
        ApiClient client = this.getKubeClient();
        return new CustomObjectsApi(client);
    }

    public CoreV1Api getCoreV1Api() {
        ApiClient client = this.getKubeClient();
        return new CoreV1Api(client);
    }

    public AppsV1Api getAppsV1Api() {
        ApiClient client = this.getKubeClient();
        return new AppsV1Api(client);
    }

}
