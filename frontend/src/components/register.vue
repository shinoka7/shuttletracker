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
                //and so on need access to stops list json
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
            <select name="days" id='days'>
                <option value="mon">Monday</option>
                <option value="tues">Tuesday</option>
                <option value="wed">Wednesday</option>
                <option value="thurs">Thursday</option>
                <option value="fri">Friday</option>
            </select>
            <select name="hours" id='hours'>
                <option value="eight">8</option>
                //times
            </select> : 
            <select name="minutes" id='minutes'>
            </select>
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

$(document).ready(function() {
    var hlist = document.getElementById('hours');
    var mlist = document.getElementById('minutes');
    var start_time = 7;
    var end_time = 23;
    for(i=start_time,j=0;i<end_time;i++,j++) {
        hlist.options[j] = new Option(i,i);
    }
    for(i=0,j=0;i<60;i=i+15,j++) {
        mlist.options[j] = new Option(i,i);
    }
});

</script>
<style lang="scss">
    .container{
        margin: 20px;
    }
</style>