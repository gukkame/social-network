import { createRouter, createWebHistory } from 'vue-router'

const router = createRouter({
  history: createWebHistory(import.meta.env.BASE_URL),
  routes: [
    {
      path: '/',
      name: 'home',

      component: () => import("../pages/Home.vue")
    },
    {
      path: '/login',
      name: 'login',
     
      component: () => import('../pages/Login.vue')
    },
    {
      path: '/signup',
      name: 'signup',
     
      component: () => import('../pages/Signup.vue')
    },
    {
      path: '/profile/:id',
      name: 'profile',
     
      component: () => import('../pages/Profile.vue')
    },
    {
      path: '/profile/:id/activity',
      name: 'activity',
     
      component: () => import('../pages/Activity.vue')
    },
    {
      path: '/profile/:id/followers',
      name: 'followers',
     
      component: () => import('../pages/Following.vue')
    },
    //GROUPS
    {
      path: '/groups',
      name: 'groups',
     
      component: () => import('../pages/Groups.vue')
    },
    {
      path: '/groups/:id',
      name: 'group',
      redirect: {name: 'groupposts'},
     
      component: () => import('../pages/group/SingleGroup.vue'),
      children: [
        {
          path: 'posts',
          name: 'groupposts',
          component: () => import('../pages/group/GroupPosts.vue'),
        },
        {
          path: 'events',
          name: 'events',
          component: () => import('../pages/group/GroupEvents.vue'),
        },
        {
          path: 'admin',
          name: 'admin',
          component: () => import('../pages/group/GroupAdmin.vue'),
        },
      ],
    },
    {
      path: '/groups/:id/posts/:postid',
      name: 'onegrouppost',
      component: () => import('../pages/group/ShowGroupPost.vue')
    },

    //CATEGORIES
    {
      path: '/Go',
      name: 'gopage',
     
      component: () => import('../pages/categories/Go.vue')
    },
    {
      path: '/HTML5',
      name: 'htmlpage',
     
      component: () => import('../pages/categories/Html.vue')
    },
    {
      path: '/CSS',
      name: 'csspage',
     
      component: () => import('../pages/categories/Css.vue')
    },
    {
      path: '/JavaScript',
      name: 'jspage',
     
      component: () => import('../pages/categories/Javascript.vue')
    },
    {
      path: '/Vue.js',
      name: 'vuepage',
     
      component: () => import('../pages/categories/Vuejs.vue')
    },

    //POSTS
    {
      path: '/Go/:id',
      name: 'Go',
     
      component: () => import('../pages/ShowPost.vue')
    },
    {
      path: '/HTML5/:id',
      name: 'HTML5',
     
      component: () => import('../pages/ShowPost.vue')
    },
    {
      path: '/CSS/:id',
      name: 'CSS',
     
      component: () => import('../pages/ShowPost.vue')
    },
    {
      path: '/JavaScript/:id',
      name: 'JavaScript',
     
      component: () => import('../pages/ShowPost.vue')
    },
    {
      path: '/Vue.js/:id',
      name: 'Vue.js',
     
      component: () => import('../pages/ShowPost.vue')
    },
    
    //404 PAGE
    {
      path: '/:pathMatch(.*)*',
      name: "PageNotFound",

      component: () => import("../pages/PageNotFound.vue")
    }
  ]
})

export default router
