<template>
    <form action="/registered.php">
    <div class="container">
        <router-link to="/">Back</router-link>
        <h2>Notification Registration</h2>
        <div>
            <lable for="stops">Target Stop: </label>
            <select name="stops">
                <option value="union">Union</option>
                <option value="blitman">Blitman</option>
                //and so on
            </select>
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
        <div>
            <label for="schedule">When to be Notified: </label>
            <select name="days">
                <option value="mon">Monday</option>
                <option value="tues">Tuesday</option>
                <option value="wed">Wednesday</option>
                <option value="thurs">Thursday</option>
                <option value="fri">Friday</option>
            </select>
            <script>
                var time = {};
                for (i=7;i<23;i++) {
                    time[i] = [];
                    for(j=15;i<60;j++) {
                        time[i].push(j); 
                    }
                }
            </script>
        </div>
    </div>
    <div>
        <input type="submit" class="submit" value="Submit" />
    </div>
    </form>
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