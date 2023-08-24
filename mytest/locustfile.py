from locust import task, run_single_user
from locust import FastHttpUser


class hongshu(FastHttpUser):
    host = "http://spider-press-test.xiaohongshu.com"
    default_headers = {
        "sec-ch-ua": '"Not.A/Brand";v="8", "Chromium";v="114", "Google Chrome";v="114"',
        "sec-ch-ua-mobile": "?0",
        "sec-ch-ua-platform": '"macOS"',
    }

    @task
    def t(self):
        with self.client.request(
            "POST",
            "/api/data",
            headers={
                "accept": "application/json, text/plain, */*",
                "accept-encoding": "gzip, deflate, br",
                "accept-language": "zh-CN,zh;q=0.9",
                "batch": "true",
                "biz-type": "tq_logserver_test",
                "content-type": "application/json;charset=UTF-8",
                "origin": "https://www.xiaohongshu.com",
                "purpose": "prefetch",
                "referer": "https://www.xiaohongshu.com/",
                "sec-fetch-dest": "empty",
                "sec-fetch-mode": "cors",
                "sec-fetch-site": "same-site",
                "sec-purpose": "prefetch;prerender",
                "user-agent": "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/114.0.0.0 Safari/537.36",
            },
            data='[{"clientTime":1690164947624,"context_nameTracker":"wapT","context_platform":"PC","context_appVersion":"discovery-undefined","context_deviceModel":"","context_deviceId":"40848649ee2761868f1f35c644c3def3","context_networkType":"unknow","context_matchedPath":"/explore","context_route":"https://www.xiaohongshu.com/explore","context_userAgent":"Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/114.0.0.0 Safari/537.36","context_artifactName":"xhs-pc-web","context_artifactVersion":"2.17.8","context_userId":"64056116000000000b01432c","measurement_name":"httpRequestTiming","measurement_data":{"method":"GET","matchedPath":"/api/sns/web/v2/user/me","traceId":"c5784af05dd05ef9","status":200,"url":"//edith.xiaohongshu.com/api/sns/web/v2/user/me","duration":570.1000000238419}},{"clientTime":1690164947698,"context_nameTracker":"wapT","context_platform":"PC","context_appVersion":"discovery-undefined","context_deviceModel":"","context_deviceId":"40848649ee2761868f1f35c644c3def3","context_networkType":"unknow","context_matchedPath":"/explore","context_route":"https://www.xiaohongshu.com/explore","context_userAgent":"Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/114.0.0.0 Safari/537.36","context_artifactName":"xhs-pc-web","context_artifactVersion":"2.17.8","context_userId":"64056116000000000b01432c","measurement_name":"httpRequestTiming","measurement_data":{"method":"GET","matchedPath":"/api/sns/web/global/config","traceId":"03eca7217f6b352d","status":200,"url":"//edith.xiaohongshu.com/api/sns/web/global/config","duration":298.80000001192093}},{"clientTime":1690164947775,"context_nameTracker":"wapT","context_platform":"PC","context_appVersion":"discovery-undefined","context_deviceModel":"","context_deviceId":"40848649ee2761868f1f35c644c3def3","context_networkType":"unknow","context_matchedPath":"/explore","context_route":"https://www.xiaohongshu.com/explore","context_userAgent":"Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/114.0.0.0 Safari/537.36","context_artifactName":"xhs-pc-web","context_artifactVersion":"2.17.8","context_userId":"64056116000000000b01432c","measurement_name":"httpRequestTiming","measurement_data":{"method":"GET","matchedPath":"/api/sns/web/v2/user/me","traceId":"3349f37335f018e0","status":200,"url":"//edith.xiaohongshu.com/api/sns/web/v2/user/me","duration":720}},{"clientTime":1690164947775,"context_nameTracker":"wapT","context_platform":"PC","context_appVersion":"discovery-undefined","context_deviceModel":"","context_deviceId":"40848649ee2761868f1f35c644c3def3","context_networkType":"unknow","context_matchedPath":"/explore","context_route":"https://www.xiaohongshu.com/explore","context_userAgent":"Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/114.0.0.0 Safari/537.36","context_artifactName":"xhs-pc-web","context_artifactVersion":"2.17.8","context_userId":"64056116000000000b01432c","measurement_name":"httpRequestTiming","measurement_data":{"method":"POST","matchedPath":"/api/sns/web/v1/racing/report","traceId":"fcb420aa8a13f5bf","status":200,"url":"//edith.xiaohongshu.com/api/sns/web/v1/racing/report","duration":376.09999999403954}},{"clientTime":1690164947965,"context_nameTracker":"wapT","context_platform":"PC","context_appVersion":"discovery-undefined","context_deviceModel":"","context_deviceId":"40848649ee2761868f1f35c644c3def3","context_networkType":"unknow","context_matchedPath":"/explore","context_route":"https://www.xiaohongshu.com/explore","context_userAgent":"Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/114.0.0.0 Safari/537.36","context_artifactName":"xhs-pc-web","context_artifactVersion":"2.17.8","context_userId":"64056116000000000b01432c","measurement_name":"httpRequestTiming","measurement_data":{"method":"POST","matchedPath":"//www.xiaohongshu.com/api/sec/v1/shield/webprofile","traceId":"9226228f11e1ce6e","status":200,"url":"//www.xiaohongshu.com/api/sec/v1/shield/webprofile","duration":319.5}}]',
            catch_response=True,
        ) as resp:
            pass

if __name__ == "__main__":
    run_single_user(hongshu)
