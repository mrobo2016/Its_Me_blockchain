import Vue from 'vue'
import Router from 'vue-router'
import Customers from './views/Customers.vue'
import CustomerDetails from './views/CustomerDetails.vue'
import Professor from './views/Professor.vue'
import ProfessorClass from './views/ProfessorClass.vue'

Vue.use(Router)

const router =  new Router({
  routes: [
    {
      path: '/',
      redirect: '/customers'
    },
    {
      path: '/customers',
      name: 'customers',
      component: Customers
    },
    {
      path: '/customerdetails/:id',
      name: 'customerdetails',
      component: CustomerDetails
    },
    {
      path: '/professor',
      name: 'professor',
      component: Professor
    },
    {
      path: '/professorClass/:id',
      name: 'professorClass',
      component: ProfessorClass
    }
  ]
})
export default router
