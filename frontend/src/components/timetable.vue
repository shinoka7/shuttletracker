<template>
	<div id='timetable'></div>
</template>

<script lang="ts">
import Vue from 'vue'
import EventBus from '../event_bus.ts'
const Tabulator = require('tabulator-tables')
export default Vue.extend({
	data () {
		return {
			tabulator: null,
			tableColumn: [
				{title: 'Route', field: 'route', align: 'center', headerSort: false, cellClick: this.notify},
			],
		}
	},

	methods: {
		
	},

	mounted() {
		this.tabulator = new Tabulator('#timetable', {
			data: [],
			columns: this.tableColumn,
			height: 250,
			layout: 'fitColumns',
			placeholder: 'Tracking Shuttles...',
		});
	},

	watch: {
		tableData: {
			handler: function(newData) {
				this.tabulator.replaceData(newData)
			},
			deep: true,
		},
	},
});
</script>

<style>
#timetable {
  background-color:#666;
  border: 1px solid #333;
  border-radius: 10px;
}

#timetable .tabulator-header {
  background-color:#666;
  color:#fff;
}

#timetable .tabulator-header .tabulator-col,
#timetable .tabulator-header .tabulator-col-row-handle {
  white-space: normal;
  background-color:#333;
}

#timetable .tabulator-tableHolder .tabulator-table .tabulator-row {
  background-color:#666;
  color:#fff;
}

#timetable .tabulator-tableHolder .tabulator-table .tabulator-row:nth-child(even) {
  background-color:#444;
}
</style>