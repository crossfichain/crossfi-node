<template lang="pug">
  div
    .search__container
      .search__container__right
        tm-select-version
        .search(@click="$emit('search', true)" v-if="$themeConfig.algolia")
          .search__icon
            icon-search
          .search__text Search
    .h1 {{$frontmatter.title}}
    .intro
      .p {{$frontmatter.description}}
    tm-help-support
</template>

<style lang="stylus" scoped>
/deep/
  .container h1
    margin-bottom 1.5rem

.search
  display flex
  align-items center
  color rgba(22, 25, 49, 0.65)
  //- width calc(var(--aside-width) - 6rem)
  cursor pointer
  transition color .15s ease-out

  &:hover
    color var(--color-text, black)

  &__container
    display flex
    justify-content flex-end
    margin-top 1rem
    margin-bottom 1rem
    padding-top 1rem

    &__right
      display grid
      grid-auto-flow column
      gap 2rem

  &__icon
    width 1.5rem
    height 1.5rem
    fill #aaa
    margin-right 0.5rem
    transition fill .15s ease-out

  &:hover &__icon
    fill var(--color-text, black)

.intro
  width 100%
  max-width 800px

.h1
  font-size 3rem
  font-weight 700
  margin-bottom 1.5rem
  line-height 3.25rem
  letter-spacing -0.02em
  padding-top 2.5rem

.h2
  font-size 2rem
  font-weight 700
  margin-top 4.5rem
  margin-bottom 1rem
  line-height 2.25rem
  letter-spacing -0.01em

.p
  font-size 1.5rem
  line-height 2.25rem

  &__alt
    margin-top 0.75rem
    margin-bottom 2rem
    font-size 1.25rem
    line-height 1.75rem

.features
  display grid
  grid-template-columns repeat(auto-fit, minmax(300px, 1fr))
  gap 1.5rem
  margin-bottom 5rem
  margin-top 2.25rem

  &__item
    cursor pointer
    display grid
    grid-auto-flow column
    grid-template-columns minmax(6rem, 1fr) 2fr
    box-shadow 0px 2px 4px rgba(22, 25, 49, 0.05), 0px 0px 1px rgba(22, 25, 49, 0.2), 0px 0.5px 0px rgba(22, 25, 49, 0.05)
    position relative
    border-radius .5rem
    background linear-gradient(302.07deg, #FFFFFF 48.96%, #EBEDFF 100%)
    outline none
    transition box-shadow 0.25s ease-out, transform 0.25s ease-out, opacity 0.4s ease-out

    &:hover:not(:active),
    &:focus:not(:active)
      box-shadow 0px 12px 24px rgba(22, 25, 49, 0.07), 0px 4px 8px rgba(22, 25, 49, 0.05), 0px 1px 0px rgba(22, 25, 49, 0.05)
      transform translateY(-2px)
      transition-duration 0.1s

    &:active
      opacity 0.7
      transition-duration 0s

    &__dark
      background linear-gradient(112.22deg, #161831 0%, #2E3148 100%)

    &__dark &__text__h2
      color white
      opacity .5

    &__dark &__text__h1
      color white

    &__dark &__text__p
      color white
      opacity .8

    &__icon
      position absolute
      top 0
      right 0
      padding .75rem
      width 1rem
      height 1rem
      fill white
      opacity .35

    &:hover &__icon
      opacity .6

    &__image
      display flex
      align-items center
      justify-content center
      align-self center
      max-height 10rem
      transition transform 0.2s ease-out

      &__img
        max-height 14rem
        max-width 10rem
        min-width 8rem

    &:hover:not(:active) &__image
      transform translateY(-0.25rem) scale(1.02)
      transition-duration 0.1s

    &__text
      padding 1.75rem 2rem 2rem
      display flex
      flex-direction column
      justify-content center

      &__h2
        font-size .75rem
        letter-spacing 0.2em
        text-transform uppercase
        color var(--color-text-dim, inherit)
        margin-bottom .25rem

      &__h1
        font-size 1.25rem
        color var(--color-text, black)
        line-height 1.75rem
        letter-spacing .01em
        font-weight 600
        margin-top .25rem
        margin-bottom .75rem

      &__p
        color var(--color-text-dim, inherit)
        font-size .875rem
        line-height 1.25rem

.sections
  display grid
  margin-top 3rem
  margin-bottom 5rem
  grid-template-columns repeat(auto-fit, minmax(300px, 1fr))
  gap 1.5rem

  &__item
    position relative
    color initial
    border-radius 0.5rem
    padding 1.5rem 1.5rem 1.5rem 5.5rem
    box-shadow 0px 2px 4px rgba(22, 25, 49, 0.05), 0px 0px 1px rgba(22, 25, 49, 0.2), 0px 0.5px 0px rgba(22, 25, 49, 0.05)
    transition box-shadow 0.25s ease-out, transform 0.25s ease-out, opacity 0.4s ease-out

    &:hover:not(:active)
      box-shadow 0px 12px 24px rgba(22, 25, 49, 0.07), 0px 4px 8px rgba(22, 25, 49, 0.05), 0px 1px 0px rgba(22, 25, 49, 0.05)
      transform translateY(-2px)
      transition-duration 0.1s

    &:active
      transition-duration 0s
      opacity 0.7

    &__icon
      position absolute
      top 1rem
      left 1.25rem
      font-size 1.5rem
      display flex
      align-items center
      justify-content center
      width 3rem
      height 3rem

    &__title
      font-weight 600
      margin-bottom 0.5rem

    &__desc
      font-size 0.875rem
      line-height 1.25rem
      color var(--color-text-dim, inherit)

.stack
  display grid
  gap 1.5rem
  grid-template-columns repeat(auto-fit, minmax(300px, 1fr))
  margin-bottom 4rem

  &__item
    position relative
    min-height 120px
    display flex
    align-items center
    padding 2rem 1.25rem
    border-radius 0.5rem
    box-shadow 0px 2px 4px rgba(22, 25, 49, 0.05), 0px 0px 1px rgba(22, 25, 49, 0.2), 0px 0.5px 0px rgba(22, 25, 49, 0.05)
    color var(--color-text, black)
    background white
    transition box-shadow 0.25s ease-out, transform 0.25s ease-out, opacity 0.4s ease-out

    &:hover:not(:active)
      box-shadow 0px 12px 24px rgba(22, 25, 49, 0.07), 0px 4px 8px rgba(22, 25, 49, 0.05), 0px 1px 0px rgba(22, 25, 49, 0.05)
      transform translateY(-2px)
      transition-duration 0.1s

    &:active
      opacity 0.7
      transition-duration 0s

    &__icon
      position absolute
      top 0
      right 0
      padding 1rem
      opacity .26
      transition opacity .15s

    &:hover &__icon
      opacity .36

    &__h1
      font-size 1.25rem
      line-height 1.5rem
      margin-bottom 0.75rem
      font-weight 600

    &__p
      font-size 0.875rem
      color rgba(22, 25, 49, 0.65)
      line-height 1.25rem

    &__wrapper
      display grid
      grid-auto-flow column
      grid-template-columns 30% 1fr
      gap 1.25rem

    &:before
      position absolute
      top 0
      left 0
      content ''
      width 50%
      height 100%
      background linear-gradient(to right, var(--accent), rgba(255, 255, 255, 0))
      border-radius 0.5rem
      opacity 0.1

    &__logo
      height 72px
      width auto

@media screen and (max-width: 1136px)
  .p
    font-size 1.25rem
    line-height 1.75rem

@media screen and (max-width: 832px)
  .h1
    padding-top 3.5rem

  .search__container
    display none

@media screen and (max-width: 752px)
  .search
    display none

@media screen and (max-width: 500px)
  .h1
    font-size 2rem
    line-height 2.25rem
    margin-bottom 1rem

  .h2
    font-size 1.5rem
    line-height 2rem
    margin-top 3rem
    margin-bottom 0.75rem

  .p__alt
    font-size 1rem
    line-height 1.5rem

  .features
    margin-bottom 1.5rem
    grid-template-columns repeat(auto-fit, minmax(240px, 1fr))

    &__item
      display block

      &:not(:active)
        box-shadow 0px 24px 40px rgba(0, 0, 0, 0.1), 0px 10px 16px rgba(0, 0, 0, 0.08), 0px 1px 0px rgba(0, 0, 0, 0.05)

      &__image
        max-height 9rem
        padding-top 1rem

      &__text
        padding 1.5rem

  .sections
    gap 0
    margin-bottom 0
    margin-top 2rem
    grid-template-columns repeat(auto-fit, minmax(240px, 1fr))
    margin-left -1rem
    margin-right -1rem

    &__item
      margin-bottom 0
      padding 1.25rem 1rem 0 5.5rem

      &,
      &:hover:not(:active)
        box-shadow none

      &__icon
        top 1rem
        left 1.25rem

      &__wrapper
        padding-bottom 1.25rem
        border-bottom 1px solid rgba(140, 145, 177, 0.32)

      &:last-child .sections__item__wrapper
        border-bottom none

    &__wrapper
      position relative
      padding 0.1px 1rem 1rem
      background white
      border-radius 0.5rem

      &:before
        position absolute
        content ''
        top 0
        left 0
        right 0
        bottom 0
        border-radius 0.5rem
        box-shadow 0px 24px 40px rgba(0, 0, 0, 0.1), 0px 10px 16px rgba(0, 0, 0, 0.08), 0px 1px 0px rgba(0, 0, 0, 0.05)

  .stack
    gap 0.75rem
    grid-template-columns repeat(auto-fill, minmax(240px, 1fr))
    margin-bottom 3rem

    &__item
      padding 1.25rem

      &:not(:active)
        box-shadow 0px 24px 40px rgba(22, 25, 49, 0.1), 0px 10px 16px rgba(22, 25, 49, 0.08), 0px 1px 0px rgba(22, 25, 49, 0.05)

      &__wrapper
        grid-template-columns 3rem 1fr

      &__h1
        font-size inherit
        line-height inherit
        margin-bottom 0.5rem

</style>
