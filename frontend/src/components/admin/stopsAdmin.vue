<template>
    <div style="margin-top: 50px;" class="container">
        <table class="table">
            <thead>
                <tr>
                    <th><abbr title="Name">Name</abbr></th>
                    <th><abbr title="ID">ID</abbr></th>
                    <th><abbr title="Desc">Description</abbr></th>
                    <th><abbr title="Lat.">Latitude</abbr></th>
                    <th><abbr title="Lng.">Longitude</abbr></th>
                    <th><abbr title="Enabled">Enabled</abbr></th>
                    <th></th>
                    <th></th>
                </tr>
                <tr v-for="stop in stops" :key="stop.id">
                    <th>{{stop.name}}</th>
                    <th></th>
                    <th></th>
                    <th></th>
                    <th></th>
                    <th></th>
                    <th></th>
                    <th></th>
                </tr>
            </thead>
            <tbody>
                
            </tbody>
        </table>
    </div>
</template>
<script lang="ts">
import Vue from 'vue';
import Stop from '@/structures/stop';
export default Vue.extend({
    name: 'stops',

    computed: {
        stops(): Stop[] {
            return this.$store.state.Stops;
        },
    },

});
</script>
