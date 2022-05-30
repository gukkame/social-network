<template>
  <div class="signup">
    <div class="formIcon">
      <img class="formImg" src="../assets/images/logo.svg" />
    </div>
    <div class="form">
      <Form @submit="signup" v-slot="{ errors }" :validation-schema="errorSchema">
        <div class="container">
          <div class="account-details la row">
            <div class="col">
              <div>
                Username
                <i class="formstar">*</i>
              </div>
              <Field class="form-input" :class="{ forminputerror: errors.username }" as="input" name="username"
                placeholder="Example101" />
              <br />
              <span class="formErrors">{{ errors.username }}</span>
            </div>
            <div class="col">
              <div>
                Email
                <i class="formstar">*</i>
              </div>
              <Field class="form-input" :class="{ forminputerror: errors.email }" as="input" name="email"
                placeholder="example@gmail.com" />
              <br />
              <span class="formErrors">{{ errors.email }}</span>
            </div>
          </div>

          <div class="user-details la row">
            <div class="col">
              <div>
                First name
                <i class="formstar">*</i>
              </div>
              <Field class="form-input" :class="{ forminputerror: errors.firstname }" as="input" name="firstname"
                placeholder="*" />
              <br />
              <span class="formErrors">{{ errors.firstname }}</span>
            </div>

            <div class="col">
              <div>
                Last name
                <i class="formstar">*</i>
              </div>
              <Field class="form-input" :class="{ forminputerror: errors.lastname }" as="input" name="lastname"
                placeholder="*" />
              <br />
              <span class="formErrors">{{ errors.lastname }}</span>
            </div>
          </div>

          <div class="user-private-details la row">
            <div class="col">
              <div>
                Age
                <i class="formstar">*</i>
              </div>
              <Field class="form-input" :class="{ forminputerror: errors.age }" type="date" name="age"
                data-date-inline-picker="true" />

              <br />
              <span class="formErrors">{{ errors.age }}</span>
            </div>
            <div class="col">
              <div>Gender</div>
              <Field class="form-input" v-model="selected" :class="{ forminputerror: errors.gender }" as="select"
                data-date-inline-picker="true" name="gender" id="age">
                <option value="Male">Male</option>
                <option value="Female">Female</option>
                <option value="Other">Other</option>
                <option value="Unspecified">Unspecified</option>
              </Field>
              <br />
              <span class="formErrors">{{ errors.gender }}</span>
            </div>
          </div>

          <div class="passwords la row">
            <div class="col">
              <div>
                Password
                <i class="formstar">*</i>
              </div>
              <Field class="form-input" :class="{ forminputerror: errors.password }" as="input" name="password"
                type="password" placeholder="Example123" />
              <br />
              <span class="formErrors">{{ errors.password }}</span>
            </div>
            <div class="col">
              <div>
                Confirm password
                <i class="formstar">*</i>
              </div>
              <Field class="form-input" :class="{ forminputerror: errors.passwordConfirmation }" as="input"
                name="passwordConfirmation" type="password" placeholder="Example123" />
              <br />
              <span class="formErrors">{{ errors.passwordConfirmation }}</span>
            </div>
          </div>
          <span class="formErrors">{{ errormsg }}</span>
          <button class="nav-button-2 nav-items-2 s">Sign Up</button>
        </div>
      </Form>
      <RouterLink class="formDirectText" to="/login">
        <i>Already signed up? Log In</i>
      </RouterLink>
    </div>
  </div>
</template>

<script>
import { Form, Field } from "vee-validate";
import * as yup from "yup";
import axios from "axios";
import router from "../router";

export default {
  name: "vaidationschema",
  data() {
    //YUP VALIDATION
    const errorSchema = yup.object().shape({
      checkUsername: yup.boolean(),
      username: yup.string() //BACK END VALIDATION TO ADD
        .required("Required")
        .matches(/^(?=.*[a-z])/, 'Username must contain lowercase characters')
        .matches(/^(?=.*[0-9])/, 'Username must contain one number')
        .min(4, "Username must be between 4 to 20 characters")
        .max(20, "Username must be between 4 to 20 characters")
        .test("Unique username", "Username already in use",
          async function (value) {
            const payload = {
              username: `${value}`,
            }

            try {
              const res = await axios.post('http://localhost:8080/available/username', payload)
              if (res.data.value === "false") {
                return false
              }
              return true
            } catch (error) {

            }
          }),

      email: yup.string() 
        .required("Required")
        .email("Please enter a correct email address")
        .test("Unique email", "Email already in use",

          async function (value) {
            const payload = {
              email: `${value}`,
            }

            try {
              const res = await axios.post('http://localhost:8080/available/email', payload)
              if (res.data.value === "false") {
                return false
              }
              return true
            } catch (error) {
            }
          }),

      firstname: yup.string()
        .matches(/^[aA-zZ\s]+$/, "Only characters are allowed")
        .required("Required")
        .max(20, "Your first name must be less than 20 characters"),

      lastname: yup.string()
        .matches(/^[aA-zZ\s]+$/, "Only characters are allowed")
        .required("Required")
        .max(20, "Your last name must less than 20 characters"),

      age: yup.string()
        .required('Required')
        .test("age", "You must be atleast 18 years old", date => {
          const getAge = birthDate => Math.floor((new Date() - new Date(birthDate).getTime()) / 3.15576e+10)
          if (getAge(date) < 18) {
            return false
          } else {
            return true
          }
        })
        .test("age", "You are too old to Sign Up", date => {
          const getAge = birthDate => Math.floor((new Date() - new Date(birthDate).getTime()) / 3.15576e+10)
          if (getAge(date) > 100) {
            return false
          } else {
            return true
          }
        }),

      gender: yup.string() 
        .required("Required"),
      password: yup.string()
        .required('Required')
        .min(8, 'Password must be between 8 to 24 characters')
        .max(24, "Password must be between 8 to 24 characters")
        .matches(/^(?=.*[a-z])/, 'Password must contain one lowercase character')
        .matches(/^(?=.*[A-Z])/, 'Password must contain one uppercase character')
        .matches(/^(?=.*[0-9])/, 'Password must contain one number'),

      passwordConfirmation: yup.string()
        .required("Required")
        .oneOf([yup.ref('password'), null], 'Passwords must match')
    });
    return {
      selected: "Unspecified",
      errorSchema,
    };
  },
  components: {
    Form,
    Field,
  },
  methods: {
    signup(values) {
      const getAge = birthDate => Math.floor((new Date() - new Date(birthDate).getTime()) / 3.15576e+10).toString()
      values.age = getAge(values.age)
      axios.post("http://localhost:8080/signup", values)
        .then((res) => {
          if (res.data.message === "Malicious user detected") {
            return this.errormsg = res.data.message
          }
          this.errormsg = ""
          router.push("/login")
        })
        .catch((error) => { });
    },
  },

  mounted() {
    let token = document.cookie
    if (token.length != 0) {
      return router.push("/")
    }
  }
};
</script>