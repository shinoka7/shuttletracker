<template>
    <div class="container">
        <router-link to="/">Back</router-link>
        <h2>Notification Registration</h2>
        <div>
            <lable for="stop">Target Stop: </label>
                <input type="text" id="stop" name="stop" required />
        </div>
        <div>
            <label for="route">Shuttle Route: </label>
                <input type="text" id="route" name="route" required />
        <div>
            <label for="number">Phone Number: </label>
                <input type="text" id="number"
                name="number" required />
        </div>
        <div>
            <label for="carrier">Carrier Type: </label>
                <input type="text" id="carrier"
                name="carrier" required /> 
        </div>
    </div>
    <div>
        <input type="submit" class="submit" value="Submit" />
    </div>
</template>
<script lang="ts">
// This component is the register window that is displayed when the `register` button on the menu is selected

import Vue from 'vue';
export default Vue.extend({

});
</script>
<style lang="scss">
    .container{
        margin: 20px;
    }
</style>