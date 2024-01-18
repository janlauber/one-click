/// <reference path="../pb_data/types.d.ts" />
migrate((db) => {
  const dao = new Dao(db)
  const collection = dao.findCollectionByNameOrId("f44b6mj6g1tv7a1")

  collection.name = "blueprints"

  return dao.saveCollection(collection)
}, (db) => {
  const dao = new Dao(db)
  const collection = dao.findCollectionByNameOrId("f44b6mj6g1tv7a1")

  collection.name = "blueprint"

  return dao.saveCollection(collection)
})
