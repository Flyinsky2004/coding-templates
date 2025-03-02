//net/index.js
import axios from "axios";
import {message} from "ant-design-vue";

const defaultError = () => message.error('发生错误，请联系管理员。')
const defaultFailure = (message) => message.warning(message)
function getAuthToken() {
    return localStorage.getItem('authToken') || '';
}
function post(url, data, success, failure = defaultFailure, error = defaultError) {
    axios.post(url, data, {
        headers: {
            "Content-Type": "application/x-www-form-urlencoded",
            "Authorization": getAuthToken()
        },
        withCredentials: true
    }).then(({data: responseData}) => {
        if (responseData.success)
            success(responseData.message,responseData.data)
        else
            failure(responseData.message)
    }).catch(error)
}
function postJSON(url, data, success, failure = defaultFailure, error = defaultError) {
    axios.post(url, data, {
        headers: {
            "Content-Type": "application/json",
            "Authorization": getAuthToken()
        },
        withCredentials: true
    }).then(({data: responseData}) => {
        if (responseData.success)
            success(responseData.message,responseData.data)
        else
            failure(responseData.message)
    }).catch(error)
}
function get(url, data = null, success, failure = defaultFailure, error = defaultError) {
    const config = {
        withCredentials: true,
        params: data,
        headers: {
            "Authorization": getAuthToken()
        },
    };

    axios.get(url, config)
        .then(({data: responseData}) => {
            if (responseData.success)
                success(responseData.message,responseData.data)
            else
                failure(responseData.message)
        })
        .catch(error)
}


export { get, post , postJSON }