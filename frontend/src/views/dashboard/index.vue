<template>
  <v-card
      class="mx-auto"
      max-width="368"
  >
    <v-card-item title="Florida">
      <template v-slot:subtitle>
        <v-icon
            class="me-1 pb-1"
            color="error"
            icon="mdi-alert"
            size="18"
        ></v-icon>

        Extreme Weather Alert
      </template>
    </v-card-item>

    <v-card-text class="py-0">
      <v-row align="center" no-gutters>
        <v-col
            class="text-h2"
            cols="6"
        >
          64&deg;F
        </v-col>

        <v-col class="text-right" cols="6">
          <v-icon
              color="error"
              icon="mdi-weather-hurricane"
              size="88"
          ></v-icon>
        </v-col>
      </v-row>
    </v-card-text>

    <div class="d-flex py-3 justify-space-between">
      <v-list-item
          density="compact"
          prepend-icon="mdi-weather-windy"
      >
        <v-list-item-subtitle>123 km/h</v-list-item-subtitle>
      </v-list-item>

      <v-list-item
          density="compact"
          prepend-icon="mdi-weather-pouring"
      >
        <v-list-item-subtitle>48%</v-list-item-subtitle>
      </v-list-item>
    </div>

    <v-expand-transition>
      <div v-if="expand">
        <div class="py-2">
          <v-slider
              v-model="time"
              :max="6"
              :step="1"
              :ticks="labels"
              class="mx-4"
              color="primary"
              density="compact"
              show-ticks="always"
              thumb-size="10"
              hide-details
          ></v-slider>
        </div>

        <v-list class="bg-transparent">
          <v-list-item
              v-for="item in forecast"
              :key="item.day"
              :append-icon="item.icon"
              :subtitle="item.temp"
              :title="item.day"
          >
          </v-list-item>
        </v-list>
      </div>
    </v-expand-transition>

    <v-divider></v-divider>

    <v-card-actions>
      <v-btn
          :text="!expand ? 'Full Report' : 'Hide Report'"
          @click="expand = !expand"
      ></v-btn>
    </v-card-actions>
  </v-card>

  <template>
    <div>
      <v-sheet
          tile
          height="54"
          color="grey lighten-3"
          class="d-flex"
      >
        <v-btn
            icon
            class="ma-2"
            @click="$refs.calendar.prev()"
        >
          <v-icon>mdi-chevron-left</v-icon>
        </v-btn>
        <v-select
            v-model="type"
            :items="types"
            dense
            outlined
            hide-details
            class="ma-2"
            label="type"
        ></v-select>
        <v-select
            v-model="mode"
            :items="modes"
            dense
            outlined
            hide-details
            label="event-overlap-mode"
            class="ma-2"
        ></v-select>
        <v-select
            v-model="weekday"
            :items="weekdays"
            dense
            outlined
            hide-details
            label="weekdays"
            class="ma-2"
        ></v-select>
        <v-spacer></v-spacer>
        <v-btn
            icon
            class="ma-2"
            @click="$refs.calendar.next()"
        >
          <v-icon>mdi-chevron-right</v-icon>
        </v-btn>
      </v-sheet>
      <v-sheet height="600">
        <v-calendar
            ref="calendar"
            v-model="value"
            :weekdays="weekday"
            :type="type"
            :events="events"
            :event-overlap-mode="mode"
            :event-overlap-threshold="30"
            :event-color="getEventColor"
            @change="getEvents"
        ></v-calendar>
      </v-sheet>
    </div>
  </template>
</template>
<script>
export default {
  data: () => ({
    labels: { 0: 'SU', 1: 'MO', 2: 'TU', 3: 'WED', 4: 'TH', 5: 'FR', 6: 'SA' },
    expand: false,
    time: 0,
    forecast: [
      { day: 'Tuesday', icon: 'mdi-white-balance-sunny', temp: '24\xB0/12\xB0' },
      { day: 'Wednesday', icon: 'mdi-white-balance-sunny', temp: '22\xB0/14\xB0' },
      { day: 'Thursday', icon: 'mdi-cloud', temp: '25\xB0/15\xB0' },
    ],
  }),
}
</script>