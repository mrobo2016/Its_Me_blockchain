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
        <p>Professor Page</p>
        <p>김ㅇㅇ 교수님 안녕하세요!</p>
        <p>It's Me 블록체인 출결 시스템을 통해 출석을 확인하세요.</p>
        <br>
        <h5>오늘의 수업:</h5>
        <div class="card centeralign addmargin" style="width: 40rem;" v-for="customer in customerlist" :key="customer.id">
            <div class="card-body" v-on:click="setSelectedCustomer(customer.name)">
                <h5 class="card-title">{{customer.name}}</h5>
                <p class="card-text">{{customer.email}}</p>
                <p class="card-text">{{customer.phone}}</p>
                <a class="btn btn-primary" v-on:click="goToDetailsPage(customer.id)"><span style="color:white">블록체인 출석 시작하기</span></a>
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
            "url": "assets/samplejson/classlist.json"
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
            this.$router.push("/professorclass/"+id);
        }
    }
}

</script>
