/// <reference path="../pb_data/types.d.ts" />
migrate((db) => {
  const dao = new Dao(db)
  const collection = dao.findCollectionByNameOrId("nwgspeyl1n5v30f")

  collection.name = "frameworks"

  return dao.saveCollection(collection)
}, (db) => {
  const dao = new Dao(db)
  const collection = dao.findCollectionByNameOrId("nwgspeyl1n5v30f")

  collection.name = "technologies"

  return dao.saveCollection(collection)
})
