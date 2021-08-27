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
    <div class="card" v-if="customerDetails.id">
        <div class="card-header">
            수업정보
        </div>
            <div class="card-body">
                <h5 class="card-title">{{customerDetails.name}}</h5>
                <p class="card-text">수업명 : {{customerDetails.name}}</p>
                <p class="card-text">학정번호 : {{customerDetails.email}}</p>
                <p class="card-text">수업시간 : {{customerDetails.phone}}</p>
                <p class="card-text">담당교수 : {{customerDetails.city}}</p>
                <br><br>
                <div>
                    <b-form @submit="onSubmit" @reset="onReset" v-if="show">
                    <b-form-group
                        id="input-group-1"
                        label="학생 고유 HASH ID:"
                        label-for="input-1"
                        description="학생 고유 HASH ID를 입력해 주세요."
                    >
                        <b-form-input
                        id="input-1"
                        v-model="form.email"
                        placeholder="Enter HASH ID"
                        required
                        ></b-form-input>
                    </b-form-group>

                    <b-form-group id="input-group-2" label="학생 이름:" label-for="input-2">
                        <b-form-input
                        id="input-2"
                        v-model="form.name"
                        placeholder="Enter name"
                        required
                        ></b-form-input>
                    </b-form-group>

                    <b-form-group id="input-group-3" label="출석 & 지각:" label-for="input-3">
                        <b-form-select
                        id="input-3"
                        v-model="form.food"
                        :options="foods"
                        required
                        ></b-form-select>
                    </b-form-group>

                    <b-button type="submit" variant="primary">Submit</b-button>
                    <b-button type="reset" variant="danger">Reset</b-button>
                    </b-form>
                </div>
            
            <br><br>
            <a v-on:click="goToMainPage()" class="btn btn-primary"><span style="color:white">Go Back</span></a>
        </div>
    </div>


</div>

</template>

<script>

// @ is an alias to /src
import axios from 'axios'

export default {
    name: 'customerdetails',
    mounted() {
        axios({
            method: "GET",
            "url": "assets/samplejson/customer"+this.$route.params.id+".json"
        }).then(response => {
            this.customerDetails = response.data;
        }, error => {
            console.error(error);
        });
    },
    data() {
        return {
            customerDetails: {},
            form: {
                email: '',
                name: '',
                food: null,
                checked: []
            },
            foods: [{ text: 'Select One', value: null }, '출석', '지각'],
            show: true
        }
    },
    methods: {
        goToMainPage: function() {
            this.$router.push("/customers");
        },
        onSubmit(event) {
            event.preventDefault()
            alert(JSON.stringify(this.form))
        },
        onReset(event) {
            event.preventDefault()
            // Reset our form values
            this.form.email = ''
            this.form.name = ''
            this.form.food = null
            this.form.checked = []
            // Trick to reset/clear native browser form validation state
            this.show = false
            this.$nextTick(() => {
                this.show = true
            })
        }
    }
}

</script>
