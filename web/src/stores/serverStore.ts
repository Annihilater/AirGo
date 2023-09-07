import {defineStore} from "pinia";
//api
import {useSystemApi} from "/@/api/system/index";
import {ElMessage} from "element-plus";

const systemApi = useSystemApi()
export const useServerStore = defineStore("serverStore", {
    state: () => ({
        serverConfig: {
            created_at: '',
            updated_at: '',
            id: 0,
            jwt: {
                signing_key: '',
                expires_time: '',
                buffer_time: '',
                issuer: '',
            },
            system: {
                enable_register: true,
                enable_email_code: false,
                enable_login_email_code: false,
                is_multipoint: true,
                sub_name: '',
                sub_url_pre: '',
                muKey: '',
                default_goods: '',
                enabled_rebate: true,    //是否开启返利
                rebate_rate: 0,          //返利率
                enabled_deduction: true, //是否开启旧套餐抵扣
                deduction_threshold: 0,  //旧套餐抵扣阈值,大于该值则抵扣
            },
            captcha: {
                key_long: 0,
                img_width: 0,
                img_height: 0,
                open_captcha: 0,
                open_captcha_time_out: 0,
            },
            pay: {
                return_url: '',
                app_id: '',
                private_key: '',
                ali_public_key: '',
                encrypt_key: '',
            },
            email: {
                email_from: '',
                email_secret: '',
                email_host: '',
                email_port: 0,
                email_is_ssl: true,
                email_nickname: '',
                email_subject: '',
                email_content: '',
            },
            rate_limit_params: {
                ip_role_param: 0,
                visit_param: 0,
            },
        } as Server,
        publicServerConfig: {
            enable_register: true,          // 是否开启注册
            enable_email_code: false,       //是否开启注册邮箱验证码
            enable_login_email_code: false, //是否开启登录邮箱验证码
            rebate_rate:0,                  //佣金率
            sub_url_pre:'',                 //订阅前缀

        },

    }),
    actions: {
        //获取系统设置
        async getServerConfig() {
            const res = await systemApi.getServerConfig()
            if (res.code === 0) {
                this.serverConfig = res.data
                ElMessage.success(res.msg)
            }
        },
        //获取公共系统设置
        async getPublicServerConfig() {
            const res = await systemApi.getPublicServerConfig()
            if (res.code === 0) {
                this.publicServerConfig = res.data
            }
        },
        //修改系统设置
        async updateServerConfig(params?: object) {
            const res = await systemApi.updateServerConfig(params)
            if (res.code === 0) {
                ElMessage.success(res.msg)
            }
        }

    }
})