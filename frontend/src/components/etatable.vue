<template>
	<div class="container">
		<router-link to="/">Back</router-link>
		<h1 class="title">Estimated Arrival Times</h1>
		<p class="subtitle">Shuttle Stop : <span v-html="stop_name"></span></p>
		<timetable/>
	</div>
</template>

<script lang="ts">
import timetable from './timetable.vue'
import Vue from 'vue'
import EventBus from '../event_bus'
export default Vue.extend({
	components: {
		timetable,
	},

	data () {
		const url = decodeURI(window.location.href).split('?')[1].split('&')
		return {
			id: url[0].split('=')[1],
			route: url[2].split('=')[1],
			stop: url[1].split('=')[1],
		}
	},
});
</script>

<style lang="scss">
.parent {
	padding: 20px;
}
.container{
	margin: 20px;
}
</style>