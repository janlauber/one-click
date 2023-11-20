/// <reference path="../pb_data/types.d.ts" />
migrate((db) => {
  const dao = new Dao(db)
  const collection = dao.findCollectionByNameOrId("22k6ts6gvnp46mc")

  // add
  collection.schema.addField(new SchemaField({
    "system": false,
    "id": "snj8tkvx",
    "name": "manifest",
    "type": "json",
    "required": false,
    "presentable": false,
    "unique": false,
    "options": {}
  }))

  return dao.saveCollection(collection)
}, (db) => {
  const dao = new Dao(db)
  const collection = dao.findCollectionByNameOrId("22k6ts6gvnp46mc")

  // remove
  collection.schema.removeField("snj8tkvx")

  return dao.saveCollection(collection)
})
