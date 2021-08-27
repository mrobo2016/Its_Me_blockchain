<style scoped>

.addmargin {
    margin-top: 10px;
    margin-bottom: 10px;
}

.vue-logo-back {
    background-color: white;
}

</style>

<template>

<div class="home">
    <div class="vue-logo-back">
        <img src="../assets/logo.png" width="316px" height="100px">
    </div>
    <div class="col-md-6 centeralign">
        <p>홍ㅇㅇ 학생 안녕하세요!</p>
        <p>It's Me 블록체인 출결 시스템으로 출석해주세요. </p>
        <p>출석이 확인 되면 It's Me 포인트가 적립되며, </p>
        <p>이 포인트는 신용평가 점수 산출에 사용 될 수 있습니다.</p>
        <a href="http://211.218.53.149:9001/main/index2">링크를 통해 신용평가 점수를 확인해 보세요.</a>
        <br> <br> <br> 
        <h5>오늘의 강의:</h5>
        <div class="card centeralign addmargin" style="width: 40rem;" v-for="customer in customerlist" :key="customer.id">
            <div class="card-body" v-on:click="setSelectedCustomer(customer.name)">
                <h5 class="card-title">{{customer.name}}</h5>
                <p class="card-text">{{customer.email}}</p>
                <p class="card-text">{{customer.phone}}</p>
                <a class="btn btn-primary" v-on:click="goToDetailsPage(customer.id)"><span style="color:white">블록체인 출석확인</span></a>
            </div>
        </div>
    </div>
    <Display v-if="selectedCustomer!=''" :selectedCustomer="selectedCustomer" />
</div>

</template>

<script>

// @ is an alias to /src
import Display from '@/components/Display.vue'
import axios from 'axios'

export default {
    name: 'customers',
    mounted() {
        axios({
            method: "GET",
            "url": "assets/samplejson/customerlist.json"
        }).then(response => {
            this.customerlist = response.data;
        }, error => {
            // eslint-disable-next-line
            console.error(error);
        });
    },
    data() {
        return {
            customerlist: [],
            selectedCustomer: ""
        }
    },
    components: {
        Display
    },
    methods: {
        setSelectedCustomer: function(name) {
            this.selectedCustomer = name;
        },
        goToDetailsPage: function(id) {
            this.$router.push("/customerdetails/"+id);
        }
    }
}

</script>
