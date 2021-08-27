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
                <p class="card-text">학정번호 : {{customerDetails.email}}</p>
                <p class="card-text">수업시간 : {{customerDetails.phone}}</p>
                <br><br>
                <div>
                    <b-form @submit="onSubmit" @reset="onReset" v-if="show">
                    <b-form-group
                        id="input-group-1"
                        label="강의자 고유 HASH ID:"
                        label-for="input-1"
                        description="강의자 고유 HASH ID를 입력해 주세요."
                    >
                        <b-form-input
                        id="input-1"
                        v-model="form.email"
                        placeholder="Enter HASH ID"
                        required
                        ></b-form-input>
                    </b-form-group>


                    <b-form-group id="input-group-3" label="출석학생에게 부여할 It's Me 포인트:" label-for="input-3">
                        <b-form-select
                        id="input-3"
                        v-model="form.present"
                        :options="present"
                        required
                        ></b-form-select>
                    </b-form-group>

                    <b-form-group id="input-group-4" label="지각학생에게 부여할 It's Me 포인트:" label-for="input-4">
                        <b-form-select
                        id="input-4"
                        v-model="form.late"
                        :options="late"
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
            "url": "assets/samplejson/professorclass"+this.$route.params.id+".json"
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
                present: null,
                late: null,
                checked: []
            },
            present: [{ text: 'Select One', value: null }, '1 It\'s Me 포인트', '2 It\'s Me 포인트', '3 It\'s Me 포인트', '4 It\'s Me 포인트', '5 It\'s Me 포인트'],
            late: [{ text: 'Select One', value: null }, '1 It\'s Me 포인트', '2 It\'s Me 포인트', '3 It\'s Me 포인트', '4 It\'s Me 포인트', '5 It\'s Me 포인트'],
            show: true
        }
    },
    methods: {
        goToMainPage: function() {
            this.$router.push("/Professor");
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
